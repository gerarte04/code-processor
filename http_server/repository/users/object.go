package users

import (
	"http_server/repository"
	"http_server/repository/models"

	"github.com/google/uuid"
)

type UsersRepo struct {
    users map[uuid.UUID]*models.User
}

func NewUsersRepo() (*UsersRepo) {
    return &UsersRepo {
        users: make(map[uuid.UUID]*models.User),
    }
}

func (db *UsersRepo) GetUser(key uuid.UUID) (*models.User, error) {
    value, ok := db.users[key]

    if !ok {
        return nil, repository.ErrorUserNotFound
    }

    return value, nil
}

func (db *UsersRepo) GetUserByCred(login string, password string) (*models.User, error) {
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

func (db *UsersRepo) PostUser(key uuid.UUID, login string, password string) error {
    if _, ok := db.users[key]; ok {
        return repository.ErrorUserKeyAlreadyUsed
    }

    for _, v := range db.users {
        if v.Login == login {
            return repository.ErrorUserAlreadyExists
        }
    }

    db.users[key] = &models.User{
        Id: key,
        Login: login,
        Password: password,
    }

    return nil
}
