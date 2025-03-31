package auth

import (
	"http_server/repository"
	"net/http"
	"strings"
)

const (
    sessionIdHeaderName = "Authorization"
)

type AuthMiddleware struct {
    sessStg repository.SessionStorage
}

func NewObject(sessStg repository.SessionStorage) *AuthMiddleware {
    return &AuthMiddleware{
        sessStg: sessStg,
    }
}

func GetSessionId(r *http.Request) (string, error) {
    value := strings.Split(r.Header.Get(sessionIdHeaderName), " ")
    
    if len(value) != 2 || value[0] != "Bearer" {
        return "", ErrorUnauthorized
    }

    return value[1], nil
}

func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if sessId, err := GetSessionId(r); err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        } else if _, err := am.sessStg.GetSession(sessId); err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}
