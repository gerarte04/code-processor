package tasks

import (
	"fmt"
	"http_server/config"
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
