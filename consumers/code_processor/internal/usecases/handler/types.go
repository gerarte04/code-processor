package handler

import (
	"code_processor/internal/models"
	"code_processor/internal/usecases"
)

type ResponseObject struct {
	TaskId string `json:"task_id"`
	Output *string `json:"output"`
	StatusCode int64 `json:"status_code"`
}

func CreateResponseObject(code *models.Code, resp *usecases.ProcessingServiceResponse) *ResponseObject {
	return &ResponseObject{
		TaskId: code.TaskId,
		Output: resp.Output,
		StatusCode: resp.StatusCode,
	}
}
