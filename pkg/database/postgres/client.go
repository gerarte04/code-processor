package postgres

import (
	"cpapp/pkg/config/types"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresClient(cfg types.PostgreSQLConfig) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
        cfg.Host, cfg.Port, cfg.DB, cfg.User, cfg.Password,
    )
    db, err := sqlx.Connect("postgres", connStr)

    if err != nil {
        return nil, err
    }

	return db, nil
}
