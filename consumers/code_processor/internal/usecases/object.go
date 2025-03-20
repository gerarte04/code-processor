package usecases

import (
	"code_processor/internal/api"
	"code_processor/internal/models"
)

type ProcessingService interface {
    Process(code *models.Code) (*api.ProcessingServiceResponse, error)
}

type ResponseWriter interface {
    WriteResponse(resp any) error
    WriteError(taskId string, err error)
}
