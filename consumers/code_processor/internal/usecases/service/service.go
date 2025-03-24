package service

import (
	"code_processor/internal/models"
	"code_processor/internal/repository"
	"code_processor/internal/usecases"
	"errors"
	"fmt"
	"strings"
)

type TasksService struct {
    procService usecases.ProcessingService
    tasksRepo repository.TasksRepo
}

func NewTasksService(procService usecases.ProcessingService, tasksRepo repository.TasksRepo) *TasksService {
    return &TasksService{procService: procService, tasksRepo: tasksRepo}
}

func (s *TasksService) ServeError(errMsg string, task *models.Task) error {
    task.Output = errMsg
    task.StatusCode = -1
    _ = s.tasksRepo.PutResult(task.Id, task)
    
    return errors.New(errMsg)
}

func (s *TasksService) ServeTask(task *models.Task) error {
    resp, err := s.procService.Process(task)

    if err != nil {
        return s.ServeError(fmt.Sprintf("processing task: %s", err.Error()), task)
    }

    task.Output = strings.ReplaceAll(resp.Output, "\u0000", "")
    task.StatusCode = resp.StatusCode
    err = s.tasksRepo.PutResult(task.Id, task)

    if err != nil {
        return s.ServeError(fmt.Sprintf("writing response: %s", err.Error()), task)
    }

    return nil
}
