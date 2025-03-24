package models

import "github.com/google/uuid"

type Task struct {
    Id uuid.UUID `db:"id"`
    Finished bool `db:"finished"`

    Output string `db:"output"`
    StatusCode int64 `db:"status_code"`

    Translator string `db:"translator"`
    Code string `db:"code"`
}
