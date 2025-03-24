package rabbitmq

import "http_server/repository/models"

type RabbitMQMessage struct {
    TaskId string `json:"task_id"`
    Translator string `json:"translator"`
    Code *string `json:"code"`
}

func CreateRabbitMQMessage(task *models.Task) *RabbitMQMessage {
    return &RabbitMQMessage{
        TaskId: task.Id.String(),
        Translator: task.Translator,
        Code: &task.Code,
    }
}
