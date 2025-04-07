package users

import (
	"cpapp/http_server/repository"

	"github.com/google/uuid"
)

type UsersService struct {
    usersRepo repository.UsersRepo
    sessStg repository.SessionStorage
}

func NewObject(usersRepo repository.UsersRepo, sessStg repository.SessionStorage) *UsersService {
    return &UsersService{
        usersRepo: usersRepo,
        sessStg: sessStg,
    }
}

func (rs *UsersService) RegisterUser(login string, password string) error {
    return rs.usersRepo.PostUser(uuid.New(), login, password)
}

func (rs *UsersService) LoginUser(login string, password string) (string, error) {
    if user, err := rs.usersRepo.GetUserByCred(login, password); err != nil {
        return "", err
    } else if sess, err := rs.sessStg.CreateSession(user.Id); err != nil {
        return "", err
    } else {
        return sess.SessionId, nil
    }
}
