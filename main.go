package main

import (
	"flag"
	"fmt"
	"http_server/api/http"
	pkgHttp "http_server/pkg/http"
	"http_server/repository/database"
	"http_server/usecases/service"

	"github.com/go-chi/chi/v5"
)

func main() {
    addr := flag.String("port", ":8080", "specify listening port")
    flag.Parse()

    db := database.NewDatabase()
    service := service.NewObject(db)
    handler := http.NewHandler(service)
    
    r := chi.NewRouter()
    handler.WithObjectHandlers(r)

    err := pkgHttp.CreateServer(*addr, r)

    if err != nil {
        _ = fmt.Errorf("%s", "failed to start: " + err.Error())
    }
}
