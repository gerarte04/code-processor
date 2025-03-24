package tasks

import (
	"http_server/repository"
	"http_server/repository/models"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
    PqUniqueViolation = "23505"
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

func (r *TasksRepo) GetTask(key uuid.UUID) (*models.Task, error) {
    row := r.db.QueryRowx("SELECT * FROM tasks WHERE id = $1", key)
    var task models.Task
    
    if err := row.StructScan(&task); err != nil {
        log.Printf("getting task: %s", err.Error())
        return nil, repository.ErrorInternalQueryError
    }

    return &task, nil
}

func (r *TasksRepo) PostTask(key uuid.UUID, task *models.Task) error {
    res, err := r.db.Exec(`INSERT INTO tasks (id, finished, translator, code)
        VALUES ($1, false, $2, $3)`, key, &task.Translator, &task.Code)

        if err != nil {
            log.Printf("posting task: %s", err.Error())
    
            if err.(*pq.Error).Code == PqUniqueViolation {
                return repository.ErrorUserAlreadyExists
            } else {
                return repository.ErrorInternalQueryError
            }
        }

    if n, err := res.RowsAffected(); err != nil || n != 1 {
        return repository.ErrorTaskKeyAlreadyUsed
    }

    return nil
}

func (r *TasksRepo) PutResult(key uuid.UUID, task *models.Task) error {
    res, err := r.db.Exec(`UPDATE tasks SET finished = true, output = $1, status_code = $2`,
        &task.Output, task.StatusCode)

    if err != nil {
        log.Printf("putting result: %s", err.Error())
        return repository.ErrorInternalQueryError
    }

    if n, err := res.RowsAffected(); err != nil || n != 1 {
        return repository.ErrorTaskNotFound
    }
    
    return nil
}
