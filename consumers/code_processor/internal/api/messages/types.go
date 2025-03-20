package messages

import (
	"code_processor/internal/api"
	"code_processor/internal/models"
)

type ResponseObject struct {
	TaskId string `json:"task_id"`
	Output *string `json:"output"`
	StatusCode int64 `json:"status_code"`
}

func CreateResponseObject(code *models.Code, resp *api.ProcessingServiceResponse) *ResponseObject {
	return &ResponseObject{
		TaskId: code.TaskId,
		Output: resp.Output,
		StatusCode: resp.StatusCode,
	}
}
