package tasks

import (
	"cpapp/consumers/code_processor/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type TasksRepo struct {
    db *sqlx.DB
}

func NewTasksRepo(db *sqlx.DB) *TasksRepo {
    return &TasksRepo{
        db: db,
    }
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
