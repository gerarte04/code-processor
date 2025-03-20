package api

import "code_processor/internal/models"

type MessageHandler interface {
    HandleMessage(message *models.Code)
}

type ProcessingServiceResponse struct {
    Output *string `json:"output"`
    StatusCode int64 `json:"status_code"`
}

type ErrorDetail struct {
    TaskId string `json:"task_id"`
    Error string `json:"output"`
    StatusCode int64 `json:"status_code"`
}
