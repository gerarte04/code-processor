package messages

import (
	"code_processor/internal/models"
	"code_processor/internal/usecases"
	"encoding/json"
	"log"
)

type MessageHandler struct {
    tasksService usecases.TasksService
}

func NewMessageHandler(tasksService usecases.TasksService) *MessageHandler {
    return &MessageHandler{
        tasksService: tasksService,
    }
}

func (h *MessageHandler) HandleMessage(message []byte) {
    var code models.Task

    if err := json.Unmarshal(message, &code); err != nil {
        log.Printf("converting message to json: %s", err.Error())
        return
    }

    if err := h.tasksService.ServeTask(&code); err != nil {
        log.Printf("%s", err.Error())
    }
}
