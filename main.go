package main

import (
	"fmt"
	"http_server/api/http"
	"http_server/config"
	_ "http_server/docs"
	"http_server/middlewares/auth"
	pkgHttp "http_server/pkg/http"
	tasksRepo "http_server/repository/tasks"
	usersRepo "http_server/repository/users"
	"http_server/usecases/sessions"
	tasksService "http_server/usecases/tasks"
	usersService "http_server/usecases/users"

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
    var cfg config.HttpConfig
    config.LoadConfig(appFlags.ConfigPath, &cfg)

    tasksRepo := tasksRepo.NewTasksRepo()
    usersRepo := usersRepo.NewUsersRepo()
    sessMgr := sessions.NewSessionManager(cfg.SessionLivingTime)
    tasksService := tasksService.NewObject(tasksRepo)
    usersService := usersService.NewObject(usersRepo, sessMgr)
    handler := http.NewHandler(tasksService, usersService)
    authMiddleware := auth.NewObject(sessMgr)
    
    r := chi.NewRouter()
    r.Use(authMiddleware.Authenticate)
    r.Get("/swagger/*", httpSwagger.WrapHandler)
    handler.WithObjectHandlers(r)
    err := pkgHttp.CreateServer(cfg.Host + ":" + cfg.Port, r)

    if err != nil {
        _ = fmt.Errorf("%s", "failed to start: " + err.Error())
    }
}
