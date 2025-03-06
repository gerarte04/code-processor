package usecases

import "errors"

var (
    ErrorTaskProcessing = errors.New("Task isn't finished yet")
    ErrorUnknownQuery = errors.New("Unknown query type\n")

    ErrorUserSessionExists = errors.New("There is active session for current user\n")
    ErrorNoUserSessionExists = errors.New("There is no active session for current user\n")
    ErrorNoSessionExists = errors.New("There is no active session with such id\n")

    ErrorSessionExpired = errors.New("User session has expired, please login again\n")
)
