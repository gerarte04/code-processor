package usecases

import "github.com/google/uuid"

type Object interface {
	Get(key uuid.UUID, queryType int) (any, error)
	Post(key uuid.UUID, value string) error
}
