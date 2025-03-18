package models

import "github.com/google/uuid"

type Result struct {
	TaskId uuid.UUID
	Output string
	StatusCode int64
}
