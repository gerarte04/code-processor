package tasks_service

import (
	"http_server/repository"
	"http_server/repository/models"
	"http_server/usecases/process"
	"time"

	"github.com/google/uuid"
)

type TasksService struct {
    tasksRepo repository.TasksRepo
}

func NewObject(tasksRepo repository.TasksRepo) *TasksService {
    return &TasksService{
        tasksRepo: tasksRepo,
    }
}

func (rs *TasksService) GetTask(key uuid.UUID) (*models.Task, error) {
    task, err := rs.tasksRepo.GetTask(key)

    if err != nil {
        return nil, err
    }

    return task, nil
}

func (rs *TasksService) PostTask(dur time.Duration) (*uuid.UUID, error) {
    key := uuid.New()
    err := rs.tasksRepo.PostTask(key)

    if err != nil {
        return nil, err
    }

    tsk, _ := rs.tasksRepo.GetTask(key)
    go process.SleepAndComplete(tsk, dur)

    return &key, nil
}
