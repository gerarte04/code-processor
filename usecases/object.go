package usecases

import (
	"time"

	"github.com/google/uuid"
)

type Object interface {
    Get(key uuid.UUID, queryType int) (string, error)
    Post(dur time.Duration) (string, error)
}
