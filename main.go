package main

import (
	"flag"
	"fmt"
	"http_server/api/http"
	_ "http_server/docs"
	pkgHttp "http_server/pkg/http"
	"http_server/repository/database"
	"http_server/usecases/service"
	"http_server/usecases/sessions"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title My API
// @version 1.0
// @description This is a sample server.

// @host localhost:8080
// @BasePath /
func main() {
    addr := flag.String("port", ":8080", "specify listening port")
    flag.Parse()

    db := database.NewDatabase()
    service := service.NewObject(db, sessions.NewSessionManager())
    handler := http.NewHandler(service)
    
    r := chi.NewRouter()
    r.Get("/swagger/*", httpSwagger.WrapHandler)
    handler.WithObjectHandlers(r)

    err := pkgHttp.CreateServer(*addr, r)

    if err != nil {
        _ = fmt.Errorf("%s", "failed to start: " + err.Error())
    }
}
