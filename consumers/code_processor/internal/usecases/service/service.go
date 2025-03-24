package service

import (
	"code_processor/internal/models"
	"code_processor/internal/repository"
	"code_processor/internal/usecases"
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

func (s *TasksService) ServeTask(task *models.Task) error {
    resp, err := s.procService.Process(task)

    if err != nil {
        return fmt.Errorf("processing task: %s", err.Error())
    }

    task.Output = strings.ReplaceAll(resp.Output, "\u0000", "")
    task.StatusCode = resp.StatusCode
    err = s.tasksRepo.PutResult(task.Id, task)

    if err != nil {
        return fmt.Errorf("writing response: %s", err.Error())
    }

    return nil
}
