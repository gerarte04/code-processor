package tasks_service

import (
	"http_server/repository"
	"http_server/repository/models"

	"github.com/google/uuid"
)

type TasksService struct {
    tasksRepo repository.TasksRepo
    sender repository.BrokerSender
}

func NewObject(tasksRepo repository.TasksRepo, sender repository.BrokerSender) *TasksService {
    return &TasksService{
        tasksRepo: tasksRepo,
        sender: sender,
    }
}

func (rs *TasksService) GetTask(key uuid.UUID) (*models.Task, error) {
    task, err := rs.tasksRepo.GetTask(key)

    if err != nil {
        return nil, err
    }

    return task, nil
}

func (rs *TasksService) PostTask(code *models.Code) (*uuid.UUID, error) {
    key := uuid.New()
    err := rs.tasksRepo.PostTask(key)

    if err != nil {
        return nil, err
    }

    err = rs.sender.Send(code)

    if err != nil {
        return nil, err
    }

    // tsk, _ := rs.tasksRepo.GetTask(key)
    // go process.SleepAndComplete(tsk, dur)

    return &key, nil
}
