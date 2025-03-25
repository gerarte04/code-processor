package tasks

import (
	"code_processor/config"
	"code_processor/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type TasksRepo struct {
    db *sqlx.DB
    cfg config.PostgreSQLConfig
}

func NewTasksRepo(cfg config.PostgreSQLConfig) (*TasksRepo, error) {
    connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
        cfg.Host, cfg.Port, cfg.DB, cfg.User, cfg.Password,
    )
    db, err := sqlx.Connect("postgres", connStr)

    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return &TasksRepo{
        db: db,
    }, nil
}

func (r *TasksRepo) PutResult(key string, task *models.Task) error {
    res, err := r.db.Exec(`UPDATE tasks SET finished = true, output = $1, status_code = $2
        WHERE id = $3`, task.Output, task.StatusCode, key)

    if err != nil {
        return fmt.Errorf("putting result: %s", err.Error())
    }

    if n, err := res.RowsAffected(); err != nil || n != 1 {
        return fmt.Errorf("putting result: task not found")
    }
    
    return nil
}
