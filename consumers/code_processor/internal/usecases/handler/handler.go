package handler

import (
	"code_processor/internal/models"
	"code_processor/internal/usecases"
	"log"
)

type MessageHandler struct {
    procService usecases.ProcessingService
    respWriter usecases.ResponseWriter
}

func NewMessageHandler(procService usecases.ProcessingService, respWriter usecases.ResponseWriter) *MessageHandler {
    return &MessageHandler{
        procService: procService,
        respWriter: respWriter,
    }
}

func (h *MessageHandler) HandleMessage(message *models.Code) {
    resp, err := h.procService.Process(message)

    if err != nil {
        log.Printf("processing message: %s", err.Error())
        return
    }

    err = h.respWriter.WriteResponse(CreateResponseObject(message, resp))

    if err != nil {
        log.Printf("writing response: %s", err.Error())
    }
}
