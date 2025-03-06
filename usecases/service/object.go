package service

import (
	"http_server/pkg/process"
	"http_server/repository"
	"http_server/repository/models"
	"time"

	"github.com/google/uuid"
)

type Object struct {
    repo repository.Object
}

func NewObject(repo repository.Object) *Object {
    return &Object{
        repo: repo,
    }
}

func (rs *Object) GetTask(key uuid.UUID) (*models.Task, error) {
    task, err := rs.repo.GetTask(key)

    if err != nil {
        return nil, err
    }

    return task, nil
}

func (rs *Object) PostTask(dur time.Duration) (*uuid.UUID, error) {
    key := uuid.New()
    err := rs.repo.PostTask(key)

    if err != nil {
        return nil, err
    }

    tsk, _ := rs.repo.GetTask(key)
    go process.SleepAndComplete(tsk, dur)

    return &key, nil
}
