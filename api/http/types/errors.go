package types

import "errors"

var (
    ErrorNotFoundPath = errors.New("no such url exists\n")
    ErrorInvalidKey = errors.New("error while parsing uuid\n")
)
