package repository

import (
	"main/repository/task"
	"time"

	"github.com/google/uuid"
)

type Object interface {
	Get(key uuid.UUID) (*task.Task, error)
	Post(key uuid.UUID, dur time.Duration) error
}
