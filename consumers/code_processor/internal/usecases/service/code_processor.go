package service

import (
	"code_processor/internal/models"
	"code_processor/internal/usecases"
	"errors"
)

type CodeProcessor struct {

}

func NewCodeProcessor() *CodeProcessor {
    return &CodeProcessor{}
}

func (p *CodeProcessor) Process(code *models.Code) (*usecases.ProcessingServiceResponse, error) {
    return nil, errors.New("zatychka")
}
