package repository

import "errors"

var (
	ErrorUserNotFound = errors.New("User with such key not found\n")
	ErrorTaskNotFound = errors.New("Task with such key not found\n")

	ErrorTaskKeyAlreadyUsed = errors.New("Task with such key already exists\n")
	ErrorUserKeyAlreadyUsed = errors.New("User with such key already exists\n")
	ErrorUserAlreadyExists = errors.New("User with such login already exists\n")

	ErrorWrongPassword = errors.New("Wrong password\n")
)
