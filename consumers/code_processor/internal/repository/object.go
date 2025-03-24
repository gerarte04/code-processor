package repository

import (
	"code_processor/internal/models"
)

type TasksRepo interface {
    PutResult(key string, result *models.Task) error
}
