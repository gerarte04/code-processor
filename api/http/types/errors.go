package types

import "errors"

var (
    ErrorNotFoundPath = errors.New("No such url exists\n")
    ErrorInvalidKey = errors.New("Error while parsing uuid\n")

    ErrorUnauthorized = errors.New("Unauthorized, please login\n")
)
