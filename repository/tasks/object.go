package tasks

import (
	"http_server/repository"
	"http_server/repository/models"

	"github.com/google/uuid"
)

type TasksRepo struct {
    tasks map[uuid.UUID]*models.Task
}

func NewTasksRepo() (*TasksRepo) {
	return &TasksRepo {
		tasks: make(map[uuid.UUID]*models.Task),
	}
}

func (db *TasksRepo) GetTask(key uuid.UUID) (*models.Task, error) {
	value, ok := db.tasks[key]

	if !ok {
		return nil, repository.ErrorTaskNotFound
	}

	return value, nil
}

func (db *TasksRepo) PostTask(key uuid.UUID) error {
	if _, ok := db.tasks[key]; ok {
		return repository.ErrorTaskKeyAlreadyUsed
	}

	db.tasks[key] = &models.Task{
		Id: key,
	}

	return nil
}
