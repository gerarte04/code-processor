package usecases

import "code_processor/internal/models"

type MessageHandler interface {
    HandleMessage(message *models.Code)
}

type ProcessingService interface {
    Process(code *models.Code) (*ProcessingServiceResponse, error)
}

type ResponseWriter interface {
    WriteResponse(resp any) error
    WriteError(taskId string, err error)
}

type ProcessingServiceResponse struct {
    Output *string `json:"output"`
    StatusCode int64 `json:"status_code"`
}

type ErrorDetail struct {
    TaskId string `json:"task_id"`
    Error string `json:"output"`
    Number int64 `json:"status_code"`
}
