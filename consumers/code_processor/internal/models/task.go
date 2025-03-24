package models

type Task struct {
    Id string `db:"id"`

    Output string `db:"output"`
    StatusCode int64 `db:"status_code"`

    Translator string `db:"translator" json:"translator"`
    Code string `db:"code" json:"code"`
}
