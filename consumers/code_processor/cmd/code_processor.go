package main

import (
	"code_processor/config"
	"code_processor/internal/api/messages"
	rabbMq "code_processor/internal/rabbitmq"
	"code_processor/internal/usecases/service"
	"code_processor/internal/usecases/writer"
	"log"
)

func main() {
    appFlags := config.ParseFlags()
    var cfg config.Config
    config.LoadConfig(appFlags.ConfigPath, &cfg)

    procService, err := service.NewCodeProcessor(cfg.ProcCfg)
    if err != nil {
        log.Fatalf("%s", err.Error())
        return
    }

    respWriter := writer.NewResponseWriter()
    msgHandler := messages.NewMessageHandler(procService, respWriter)
    rabbitMqReceiver, err := rabbMq.NewRabbitMQReceiver(cfg.RabbMQCfg, msgHandler)

    if err != nil {
        log.Fatalf("%s", err.Error())
    }

    if err = rabbitMqReceiver.StartReceive(); err != nil {
        rabbitMqReceiver.Close()
        log.Fatalf("%s", err.Error())
    }
}
