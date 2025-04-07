package tasks

import (
	"cpapp/http_server/repository"
	"cpapp/http_server/repository/models"
	"cpapp/pkg/database"
	"cpapp/pkg/database/postgres"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TasksRepo struct {
    db *sqlx.DB
}

func NewTasksRepo(db *sqlx.DB) *TasksRepo {
    return &TasksRepo{
        db: db,
    }
}

func (r *TasksRepo) GetTask(key uuid.UUID) (*models.Task, error) {
    row := r.db.QueryRowx("SELECT * FROM tasks WHERE id = $1", key)
    var task models.Task
    
    if err := row.StructScan(&task); err != nil {
        log.Printf("getting task: %s", err.Error())
        return nil, repository.ErrorTaskNotFound
    }

    return &task, nil
}

func (r *TasksRepo) PostTask(key uuid.UUID, task *models.Task) error {
    res, err := r.db.Exec(`INSERT INTO tasks (id, translator, code)
        VALUES ($1, $2, $3)`, key, task.Translator, task.Code)

        if err != nil {
            log.Printf("posting task: %s", err.Error())
    
            if postgres.ProcessError(err) == database.ErrorUniqueViolation {
                return repository.ErrorTaskKeyAlreadyUsed
            } else {
                return repository.ErrorInternalQueryError
            }
        }

    if n, err := res.RowsAffected(); err != nil || n != 1 {
        return repository.ErrorTaskKeyAlreadyUsed
    }

    return nil
}
