package tasks_service

import (
	"cpapp/http_server/repository"
	"cpapp/http_server/repository/models"

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

func (rs *TasksService) PostTask(task *models.Task) (*uuid.UUID, error) {
    key := uuid.New()
    err := rs.tasksRepo.PostTask(key, task)

    if err != nil {
        return nil, err
    }

    task.Id = key
    err = rs.sender.Send(task)

    if err != nil {
        return nil, err
    }

    return &key, nil
}
