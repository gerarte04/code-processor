package server

import (
	"net/http"
)

func CreateServer(addr string, handler http.Handler) error {
    s := &http.Server{
        Addr: addr,
        Handler: handler,
    }

    return s.ListenAndServe()
}
