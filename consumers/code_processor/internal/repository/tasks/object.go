package tasks

import (
	"code_processor/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type TasksRepo struct {
    db *sqlx.DB
}

func NewTasksRepo(connStr string) (*TasksRepo, error) {
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
