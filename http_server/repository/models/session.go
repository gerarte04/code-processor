package models

import (
	"github.com/google/uuid"
)

type Session struct {
    UserId uuid.UUID
    SessionId string
    // ExpiresAt time.Time
}
