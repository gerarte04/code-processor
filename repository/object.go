package repository

import (
	"http_server/repository/task"

	"github.com/google/uuid"
)

type Object interface {
	Get(key uuid.UUID) (*task.Task, error)
	Post(key uuid.UUID) error
}
