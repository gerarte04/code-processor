package types

import "errors"

var (
    ErrorNotFoundPath = errors.New("no such url exists\n")
    ErrorInvalidKey = errors.New("error while parsing uuid\n")

    ErrorUnauthorized = errors.New("Unauthorized, please login\n")
    ErrorBadCookie = errors.New("Bad cookie reading\n")
)
