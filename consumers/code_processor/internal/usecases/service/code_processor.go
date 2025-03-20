package service

import (
	"bufio"
	"code_processor/config"
	"code_processor/internal/api"
	"code_processor/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

var (
	extensions = map[string]string {
		"gcc": ".cpp",
		"clang": ".cpp",
		"python3": ".py",
	}
)

type CodeProcessor struct {
	cli *client.Client
    cfg config.ProcessorConfig
}

func NewCodeProcessor(cfg config.ProcessorConfig) (*CodeProcessor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, fmt.Errorf("creating docker client: %s", err.Error())
	}

    list, err := cli.ContainerList(context.Background(), container.ListOptions{})

    if err != nil {
        return nil, fmt.Errorf("getting containers list: %s", err.Error())
    }

    for _, item := range list {
        if slices.Contains(item.Names, "/" + cfg.ContainerName) {
            err = cli.ContainerRemove(context.Background(), item.ID, container.RemoveOptions{Force: true})

            if err != nil {
                return nil, fmt.Errorf("removing existing container: %s", err.Error())
            }
        }
    }

    return &CodeProcessor{
		cli: cli,
        cfg: cfg,
	}, nil
}

func (p *CodeProcessor) CreateCodeFile(code *models.Code) error {
    ext, ok := extensions[code.Translator]

    if !ok {
        return fmt.Errorf("no supported translators")
    }

	f, err := os.Create(fmt.Sprintf("%s/%s%s", p.cfg.ImagePath, p.cfg.CodeFileName, ext))

	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.WriteString(code.Code)

	if err != nil {
		return err
	}

	return nil
}

func (p *CodeProcessor) BuildImage(code *models.Code) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.cfg.BuildTimeout)
	defer cancel()

	tar, err := archive.TarWithOptions(p.cfg.ImagePath, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("preparing archive: %s", err.Error())
	}
	defer tar.Close()

    ext, ok := extensions[code.Translator]

    if !ok {
        return fmt.Errorf("no supported translators")
    }

	fileName := p.cfg.CodeFileName + ext

	opts := types.ImageBuildOptions{
		Dockerfile: p.cfg.Dockerfile,
		Tags: []string{p.cfg.ImageName},
		Remove:     true,
		BuildArgs: map[string]*string {
			"translator": &code.Translator,
			"file": &fileName,
		},
	}
	resp, err := p.cli.ImageBuild(ctx, tar, opts)

	if err != nil {
		return fmt.Errorf("building docker image: %s", err.Error())
	}

	defer resp.Body.Close()

	var lastLine string

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

    mp := make(map[string]string)
    _ = json.Unmarshal([]byte(lastLine), &mp)

    if d, ok := mp["error"]; ok {
        return errors.New("error while building image: " + d)
    }

	return nil
}

func (p *CodeProcessor) CreateAndRunContainer() (*api.ProcessingServiceResponse, error) {
	var resp container.CreateResponse
	resp, err := p.cli.ContainerCreate(context.Background(), &container.Config{
		Image: p.cfg.ImageName,
		Tty: false,
	}, nil, nil, nil, p.cfg.ContainerName);

    if err != nil {
        return nil, fmt.Errorf("error creating container: %s", err.Error())
    }

	defer func() {
		if err := p.cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true}); err != nil {
			log.Printf("failed to remove container: %s", err.Error())
		}
	}()

	if err := p.cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("starting docker container: %s", err.Error())
	}

    ctx, cancel := context.WithTimeout(context.Background(), p.cfg.RunTimeout)
	defer cancel()

	statusCh, errCh := p.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	var statusCode int64 = -1

	select {
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("error while waiting container: %s", err.Error())
		}
	case status := <-statusCh:
		statusCode = status.StatusCode
	}

	out, err := p.cli.ContainerLogs(context.Background(), resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})

	if err != nil {
		return nil, fmt.Errorf("reading docker container logs: %s", err.Error())
	}

	defer out.Close()

	buf, err := io.ReadAll(out)

	if err != nil {
		return nil, fmt.Errorf("reading from io.ReadCloser: %s", err.Error())
	}

	str := string(buf)

	return &api.ProcessingServiceResponse{
		Output: &str,
		StatusCode: statusCode,
	}, nil
}

func (p *CodeProcessor) Process(code *models.Code) (*api.ProcessingServiceResponse, error) {
	err := p.CreateCodeFile(code)

	if err != nil {
		return nil, err
	}

	if err = p.BuildImage(code); err != nil {
		return nil, err
	}

	var resp *api.ProcessingServiceResponse

	if resp, err = p.CreateAndRunContainer(); err != nil {
		return nil, err
	}

	return resp, nil
}
