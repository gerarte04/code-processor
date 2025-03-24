package models

type Task struct {
    Id string `db:"id" json:"task_id"`

    Output string
    StatusCode int64

    Translator string `db:"translator" json:"translator"`
    Code string `db:"code" json:"code"`
}
