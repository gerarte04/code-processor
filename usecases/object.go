package usecases

import (
	"http_server/repository/models"
	"time"

	"github.com/google/uuid"
)

type Object interface {
    GetTask(key uuid.UUID) (*models.Task, error)
    PostTask(dur time.Duration) (*uuid.UUID, error)
}
