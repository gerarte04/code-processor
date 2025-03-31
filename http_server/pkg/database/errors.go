package database

import "errors"

var (
	ErrorUniqueViolation = errors.New("Violation of key uniqueness")
	UndocumentedError = errors.New("not specific error")
)
