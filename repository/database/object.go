package database

import (
	"http_server/repository"
	"http_server/repository/models"

	"github.com/google/uuid"
)

type Object struct {
    tasks map[uuid.UUID]*models.Task
	users map[uuid.UUID]*models.User
}

func NewDatabase() (*Object) {
	return &Object {
		tasks: make(map[uuid.UUID]*models.Task),
		users: make(map[uuid.UUID]*models.User),
	}
}

func (db *Object) GetTask(key uuid.UUID) (*models.Task, error) {
	value, ok := db.tasks[key]

	if !ok {
		return nil, repository.ErrorTaskNotFound
	}

	return value, nil
}

func (db *Object) PostTask(key uuid.UUID) error {
	db.tasks[key] = &models.Task{
		Id: key,
	}

	return nil
}

func (db *Object) GetUser(key uuid.UUID) (*models.User, error) {
	value, ok := db.users[key]

	if !ok {
		return nil, repository.ErrorUserNotFound
	}

	return value, nil
}

func (db *Object) GetUserByCred(login string, password string) (*models.User, error) {
	for _, v := range db.users {
		if v.Login == login {
			if v.Password == password {
				return v, nil
			} else {
				return nil, repository.ErrorWrongPassword
			}
		}
	}

	return nil, repository.ErrorUserNotFound
}

func (db *Object) PostUser(key uuid.UUID, login string, password string) error {
	for _, v := range db.users {
		if v.Login == login {
			return repository.ErrorUserAlreadyExists
		}
	}

	db.users[key] = &models.User{
		Login: login,
		Password: password,
	}

	return nil
}
