package models

import "github.com/google/uuid"

type Code struct {
    TaskId uuid.UUID
    Translator string
    Code string
}
