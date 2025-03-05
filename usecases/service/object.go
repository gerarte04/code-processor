package service

import (
	"http_server/repository"
	"http_server/repository/task"
	"http_server/usecases"
	"strconv"
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

func (rs *Object) Get(key uuid.UUID, queryType int) (string, error) {
    task, err := rs.repo.Get(key)

    if err != nil {
        return "", err
    } else if queryType == usecases.GetResultQuery {
        if task.Finished {
            return strconv.Itoa(task.Result), nil
        } else {
            return "", usecases.ErrorTaskProcessing
        }
    } else if queryType == usecases.GetStatusQuery {
        if task.Finished {
            return "ready", nil
        } else {
            return "in_progress", nil
        }
    } else {
        return "", usecases.ErrorUnknownQuery
    }
}

func (rs *Object) Post(dur time.Duration) (string, error) {
    key := uuid.New()
    err := rs.repo.Post(key)

    if err != nil {
        return "", err
    }

    tsk, _ := rs.repo.Get(key)
    go task.SleepAndComplete(tsk, dur)

    return key.String(), nil
}
