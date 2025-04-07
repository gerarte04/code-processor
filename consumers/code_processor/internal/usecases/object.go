package usecases

import (
	"cpapp/consumers/code_processor/internal/models"
)

type ProcessingService interface {
    Process(code *models.Task) (*ProcessingServiceResponse, error)
}

type TasksService interface {
    ServeTask(task *models.Task) error
}
