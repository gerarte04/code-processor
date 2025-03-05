package database

import (
	"http_server/repository"
	"http_server/repository/task"

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
		return nil, repository.NotFound
	}

	return value, nil
}

func (db *Object) Post(key uuid.UUID) error {
	db.tasks[key] = &task.Task{
		Id: key,
	}

	return nil
}
