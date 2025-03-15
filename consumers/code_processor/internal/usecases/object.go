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
}

type ProcessingServiceResponse struct {
    Result string
    StatusCode int
}
