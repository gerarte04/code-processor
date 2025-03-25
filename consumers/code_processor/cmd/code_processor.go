package main

import (
	"code_processor/config"
	"code_processor/internal/api/messages"
	rabbMq "code_processor/internal/rabbitmq"
	tasks "code_processor/internal/repository/tasks"
	"code_processor/internal/usecases/process"
	"code_processor/internal/usecases/service"
	"log"
)

func main() {
    appFlags := config.ParseFlags()
    var cfg config.Config
    config.LoadConfig(appFlags.ConfigPath, &cfg)

    tasksRepo, err := tasks.NewTasksRepo(cfg.PostgresCfg)
    if err != nil {
        log.Fatalf("%s", err.Error())
    }

    procService, err := process.NewCodeProcessor(cfg.ProcCfg)
    if err != nil {
        log.Fatalf("%s", err.Error())
    }

    tasksService := service.NewTasksService(procService, tasksRepo)
    msgHandler := messages.NewMessageHandler(tasksService)
    rabbitMqReceiver, err := rabbMq.NewRabbitMQReceiver(cfg.RabbMQCfg, msgHandler)

    if err != nil {
        log.Fatalf("%s", err.Error())
    }

    if err = rabbitMqReceiver.StartReceive(); err != nil {
        rabbitMqReceiver.Close()
        log.Fatalf("%s", err.Error())
    }
}
