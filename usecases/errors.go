package usecases

import "errors"

var (
    ErrorTaskProcessing = errors.New("Task isn't finished yet")
    ErrorUnknownQuery = errors.New("Unknown query type\n")
)
