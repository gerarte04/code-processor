package main

import (
	rabbMq "code_processor/internal/rabbitmq"
	"code_processor/internal/usecases/handler"
	"code_processor/internal/usecases/service"
	"code_processor/internal/usecases/writer"
	"log"
)

func main() {
    procService, err := service.NewCodeProcessor()
    if err != nil {
        log.Fatalf("%s", err.Error())
        return
    }

    respWriter := writer.NewResponseWriter()
    msgHandler := handler.NewMessageHandler(procService, respWriter)
    rabbitMqReceiver, err := rabbMq.NewRabbitMQReceiver("amqp://guest:guest@broker:5672", msgHandler)

    if err != nil {
        log.Fatalf("%s", err.Error())
    }

    if err = rabbitMqReceiver.StartReceive(); err != nil {
        rabbitMqReceiver.Close()
        log.Fatalf("%s", err.Error())
    }
}
