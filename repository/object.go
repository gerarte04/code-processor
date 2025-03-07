package repository

import (
	"http_server/repository/models"

	"github.com/google/uuid"
)

type TasksRepo interface {
	GetTask(key uuid.UUID) (*models.Task, error)
	PostTask(key uuid.UUID) error
}

type UsersRepo interface {
	GetUser(key uuid.UUID) (*models.User, error)
	GetUserByCred(login string, password string) (*models.User, error)
	PostUser(key uuid.UUID, login string, password string) error
}
