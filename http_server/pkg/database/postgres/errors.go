package postgres

import (
	"http_server/pkg/database"

	"github.com/lib/pq"
)

const (
    PqUniqueViolation = "23505"
)

type PostgresErrorProcessor struct {}

func NewPostgresErrorProcessor() *PostgresErrorProcessor {
	return &PostgresErrorProcessor{}
}

func (p *PostgresErrorProcessor) ProcessError(err error) error {
	if err.(*pq.Error).Code == PqUniqueViolation {
		return database.ErrorUniqueViolation
	}

	return database.UndocumentedError
}
