package rabbitmq

import "http_server/repository/models"

type RabbitMQMessage struct {
	TaskId string `json:"task_id"`
	Translator string `json:"translator"`
	Code *string `json:"code"`
}

func CreateRabbitMQMessage(code *models.Code) *RabbitMQMessage {
	return &RabbitMQMessage{
		TaskId: code.TaskId.String(),
		Translator: code.Translator,
		Code: &code.Code,
	}
}
