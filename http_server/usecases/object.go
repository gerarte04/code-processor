package usecases

import (
	"http_server/repository/models"

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

type SessionManager interface {
    StartSession(userId uuid.UUID) (*models.Session, error)
    StopSession(sessionId string) error
    GetSession(sessionId string) (*models.Session, error)
}
