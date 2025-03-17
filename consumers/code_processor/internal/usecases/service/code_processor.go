package service

import (
	"bufio"
	"code_processor/internal/models"
	"code_processor/internal/usecases"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

var (
	extensions = map[string]string {
		"gcc": ".cpp",
		"clang": ".cpp",
		"python": ".py",
	}
)

type CodeProcessor struct {
	cli *client.Client
	ctx *context.Context
}

func NewCodeProcessor() (*CodeProcessor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, fmt.Errorf("creating docker client: %s", err.Error())
	}

    return &CodeProcessor{
		cli: cli,
	}, nil
}

func CreateCodeFile(code *models.Code) (string, error) {
	path := "./build/"
	f, err := os.Create(path + "file" + extensions[code.Translator])

	if err != nil {
		return "", err
	}

	defer f.Close()
	_, err = f.WriteString(code.Code)

	if err != nil {
		return "", err
	}

	return path, nil
}

func (p *CodeProcessor) BuildImage(path string, code *models.Code) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	tar, err := archive.TarWithOptions(path, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("preparing archive: %s", err.Error())
	}
	defer tar.Close()

	fileName := "file" + extensions[code.Translator]

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags: []string{"processing_code_image"},
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

	// var lastLine string

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		// lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	return nil
}

func (p *CodeProcessor) CreateAndRunContainer() (*usecases.ProcessingServiceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	var resp container.CreateResponse
	resp, err := p.cli.ContainerCreate(ctx, &container.Config{
		Image: "processing_code_image",
		Tty: false,
	}, nil, nil, nil, "code_container");
	tries := 3

	for err != nil {
		if tries == 0 {
			return nil, fmt.Errorf("cannot remove container: %s", err.Error())
		}

		list, err := p.cli.ContainerList(ctx, container.ListOptions{})

		for _, item := range list {
			if item.Image == "processing_code_image" {
				p.cli.ContainerRemove(ctx, item.ID, container.RemoveOptions{})
				tries -= 1

				resp, err = p.cli.ContainerCreate(ctx, &container.Config{
					Image: "processing_code_image",
					Tty: false,
				}, nil, nil, nil, "code_container");

				continue
			}
		}

		return nil, fmt.Errorf("creating docker container: %s", err.Error())
	}

	defer func() {
		if err := p.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{}); err != nil {
			log.Printf("failed to remove container: %s", err.Error())
		}
	}()

	if err := p.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("starting docker container: %s", err.Error())
	}

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

	out, err := p.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
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

	return &usecases.ProcessingServiceResponse{
		Output: &str,
		StatusCode: statusCode,
	}, nil
}

func (p *CodeProcessor) Process(code *models.Code) (*usecases.ProcessingServiceResponse, error) {
	path, err := CreateCodeFile(code)

	if err != nil {
		return nil, err
	}

	if err = p.BuildImage(path, code); err != nil {
		return nil, err
	}

	var resp *usecases.ProcessingServiceResponse

	if resp, err = p.CreateAndRunContainer(); err != nil {
		return nil, err
	}

	return resp, nil
}
