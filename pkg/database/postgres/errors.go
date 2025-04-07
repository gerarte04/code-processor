package postgres

import (
	"cpapp/pkg/database"

	"github.com/lib/pq"
)

const (
    PqUniqueViolation = "23505"
)

func ProcessError(err error) error {
	if err.(*pq.Error).Code == PqUniqueViolation {
		return database.ErrorUniqueViolation
	}

	return database.UndocumentedError
}
