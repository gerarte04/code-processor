package service

import (
	"errors"
	"main/repository"
	"main/usecases"
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

func (rs *Object) Get(key uuid.UUID, queryType int) (any, error) {
	task, err := rs.repo.Get(key)

	if err != nil {
		return nil, err
	} else if queryType == usecases.GetResultQuery {
		if task.Finished {
			return task.Result, nil
		} else {
			return nil, errors.New("task isn't finished yet")
		}
	} else if queryType == usecases.GetStatusQuery {
		return task.Finished, nil
	} else {
		return nil, errors.New("unknown query type")
	}
}

func (rs *Object) Post(key uuid.UUID, value string) error {
	dur, err := time.ParseDuration(value)

	if err != nil {
		return err
	}

	return rs.repo.Post(key, dur)
}
