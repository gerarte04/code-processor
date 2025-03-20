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

func (db *TasksRepo) PostTask(key uuid.UUID, code *models.Code) error {
    if _, ok := db.tasks[key]; ok {
        return repository.ErrorTaskKeyAlreadyUsed
    }

    newTask := models.Task{
        Id: key,
        Finished: false,
        Code: code,
    }
    newTask.Code.TaskId = key
    db.tasks[key] = &newTask

    return nil
}

func (db *TasksRepo) PutResult(key uuid.UUID, result *models.Result) error {
    if _, ok := db.tasks[key]; !ok {
        return repository.ErrorTaskNotFound
    }

    db.tasks[key].Result = result
    db.tasks[key].Finished = true
    
    return nil
}
