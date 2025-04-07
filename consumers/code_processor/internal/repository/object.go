package repository

import (
	"cpapp/consumers/code_processor/internal/models"
)

type TasksRepo interface {
    PutResult(key string, result *models.Task) error
}
