package main

import (
	"cpapp/http_server/api/http"
	"cpapp/http_server/config"
	_ "cpapp/http_server/docs"
	"cpapp/http_server/middlewares/auth"
	rabbMq "cpapp/http_server/repository/rabbitmq"
	redis "cpapp/http_server/repository/redis"
	tasksRepo "cpapp/http_server/repository/tasks"
	usersRepo "cpapp/http_server/repository/users"
	tasksService "cpapp/http_server/usecases/tasks"
	usersService "cpapp/http_server/usecases/users"
	pkgConfig "cpapp/pkg/config"
	"cpapp/pkg/database/postgres"
	pkgHttp "cpapp/pkg/http"
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title My API
// @version 1.0
// @description This is a sample server.

// @host localhost:8080
// @BasePath /cmd/
func main() {
    appFlags := config.ParseFlags()
    var cfg config.Config
    pkgConfig.LoadConfig(appFlags.ConfigPath, &cfg)

    pgDb, err := postgres.NewPostgresClient(cfg.PostgresCfg)
    if err != nil {
        log.Fatalf("%s", err.Error())
    }

    tasksRepo := tasksRepo.NewTasksRepo(pgDb)
    usersRepo := usersRepo.NewUsersRepo(pgDb)

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
