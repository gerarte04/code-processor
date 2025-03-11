package middlewares

import "net/http"

type AuthMiddleware interface {
    Authenticate(next http.Handler) http.Handler
}
