package database

import (
	"errors"
	"main/repository/task"
	"time"

	"github.com/google/uuid"
)

type Object struct {
    tasks map[uuid.UUID]*task.Task
}

func NewDatabase() (*Object) {
	return &Object {
		tasks: make(map[uuid.UUID]*task.Task),
	}
}

func (db *Object) Get(key uuid.UUID) (*task.Task, error) {
	value, ok := db.tasks[key]

	if !ok {
		return nil, errors.New("key doesn't exists")
	}

	return value, nil
}

func (db *Object) Post(key uuid.UUID, dur time.Duration) error {
	db.tasks[key] = &task.Task{
		Id: key,
	}

	go task.SleepAndComplete(db.tasks[key], dur)

	return nil
}
