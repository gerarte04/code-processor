package models

import "github.com/google/uuid"

type Task struct {
    Id uuid.UUID
    Finished bool
    Result *Result
    Code *Code
}
