package service

import (
	"bufio"
	"code_processor/internal/models"
	"code_processor/internal/usecases"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
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

	return "", nil
}

func (p *CodeProcessor) BuildImage(path string, code *models.Code) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	tar, err := archive.TarWithOptions(path, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("preparing archive: %s", err.Error())
	}

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
	res, err := p.cli.ImageBuild(ctx, tar, opts)

	if err != nil {
		return fmt.Errorf("building docker image: %s", err.Error())
	}

	defer res.Body.Close()

	// var lastLine string

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		// lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	return nil
}

func (p *CodeProcessor) Process(code *models.Code) (*usecases.ProcessingServiceResponse, error) {
	path, err := CreateCodeFile(code)

	if err != nil {
		return nil, err
	}

	err = p.BuildImage(path, code)

	if err != nil {
		return nil, err
	}

	return &usecases.ProcessingServiceResponse{
		Result: "aboba",
		StatusCode: 1,
	}, nil
}
