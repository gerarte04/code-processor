package usecases

import (
	"http_server/repository/models"
	"time"

	"github.com/google/uuid"
)

type Object interface {
    GetTask(key uuid.UUID, sessionId string) (*models.Task, error)
    PostTask(dur time.Duration, sessionId string) (*uuid.UUID, error)
	RegisterUser(login string, password string) error
	LoginUser(login string, password string) (string, error)
}

type SessionManager interface {
	StartSession(userId uuid.UUID, expiresAt time.Time) (*models.Session, error)
	StopSession(sessionId string) error
	GetSession(sessionId string) (*models.Session, error)
}
