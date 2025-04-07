package main

import (
	"cpapp/consumers/code_processor/config"
	"cpapp/consumers/code_processor/internal/api/messages"
	rabbMq "cpapp/consumers/code_processor/internal/rabbitmq"
	tasks "cpapp/consumers/code_processor/internal/repository/tasks"
	"cpapp/consumers/code_processor/internal/usecases/process"
	"cpapp/consumers/code_processor/internal/usecases/service"
	pkgConfig "cpapp/pkg/config"
	"cpapp/pkg/database/postgres"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    appFlags := pkgConfig.ParseFlags()
    var cfg config.Config
    pkgConfig.LoadConfig(appFlags.ConfigPath, &cfg)

    db, err := postgres.NewPostgresClient(cfg.PostgresCfg)
    if err != nil {
        log.Fatalf("%s", err.Error())
    }

    tasksRepo := tasks.NewTasksRepo(db)

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

    http.Handle("/metrics", promhttp.Handler())
    
    if err = http.ListenAndServe(":" + cfg.PromCfg.ListenPort, nil); err != nil {
        log.Fatalf("%s", err.Error())
    }
}
