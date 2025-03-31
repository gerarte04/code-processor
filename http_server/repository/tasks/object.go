package tasks

import (
	"http_server/pkg/database"
	"http_server/repository"
	"http_server/repository/models"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TasksRepo struct {
    db *sqlx.DB
    ep database.DBErrorProcessor
}

func NewTasksRepo(db *sqlx.DB, ep database.DBErrorProcessor) *TasksRepo {
    return &TasksRepo{
        db: db,
        ep: ep,
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
    
            if r.ep.ProcessError(err) == database.ErrorUniqueViolation {
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
