package repository

import "errors"

var (
	ErrorUserNotFound = errors.New("User with such key not found\n")
	ErrorTaskNotFound = errors.New("Task with such key not found\n")
	ErrorUserAlreadyExists = errors.New("User with such login already exists\n")
)
