package usecases

import (
	"cpapp/http_server/repository/models"

	"github.com/google/uuid"
)

type TasksService interface {
    GetTask(key uuid.UUID) (*models.Task, error)
    PostTask(task *models.Task) (*uuid.UUID, error)
}

type UsersService interface {
    RegisterUser(login string, password string) error
    LoginUser(login string, password string) (string, error)
}
