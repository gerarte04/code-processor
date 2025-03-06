package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserId uuid.UUID
	SessionId string
	ExpiresAt time.Time
}
