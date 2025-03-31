package main

import (
	"fmt"
	"http_server/api/http"
	"http_server/config"
	_ "http_server/docs"
	"http_server/middlewares/auth"
	"http_server/pkg/database/postgres"
	pkgHttp "http_server/pkg/http"
	rabbMq "http_server/repository/rabbitmq"
	redis "http_server/repository/redis"
	tasksRepo "http_server/repository/tasks"
	usersRepo "http_server/repository/users"
	tasksService "http_server/usecases/tasks"
	usersService "http_server/usecases/users"
	"log"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title My API
// @version 1.0
// @description This is a sample server.

// @host localhost:8080
// @BasePath /
func main() {
    appFlags := config.ParseFlags()
    var cfg config.Config
    config.LoadConfig(appFlags.ConfigPath, &cfg)

    pgDb, err := postgres.NewPostgresClient(cfg.PostgresCfg)
    if err != nil {
        log.Fatalf("%s", err.Error())
    }

    errProc := postgres.NewPostgresErrorProcessor()
    tasksRepo := tasksRepo.NewTasksRepo(pgDb, errProc)
    usersRepo := usersRepo.NewUsersRepo(pgDb, errProc)

    sessStg, err := redis.NewSessionStorage(cfg.ServiceCfg, cfg.RedisCfg)
    if err != nil {
        log.Fatalf("%s", err.Error())
    }

    rabbitMQSender, err := rabbMq.NewRabbitMQSender(cfg.RabbMQCfg)
    if err != nil {
        log.Fatalf("%s", err.Error())
    }

    tasksService := tasksService.NewObject(tasksRepo, rabbitMQSender)
    usersService := usersService.NewObject(usersRepo, sessStg)
    handler := http.NewHandler(tasksService, usersService)
    authMiddleware := auth.NewObject(sessStg)

    r := chi.NewRouter()
    handler.RouteHandlers(r,
        handler.WithFreeUserHandlers(r),
        handler.WithSecuredUserHandlers(r, authMiddleware),
    )
    r.Get("/swagger/*", httpSwagger.WrapHandler)

    err = pkgHttp.CreateServer(cfg.HttpCfg.Host + ":" + cfg.HttpCfg.Port, r)

    if err != nil {
        _ = fmt.Errorf("%s", "failed to start: " + err.Error())
    }
}
