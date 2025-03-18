package users

import (
	"http_server/repository"
	"http_server/usecases"

	"github.com/google/uuid"
)

type UsersService struct {
    usersRepo repository.UsersRepo
    sessMgr usecases.SessionManager
}

func NewObject(usersRepo repository.UsersRepo, sessMgr usecases.SessionManager) *UsersService {
    return &UsersService{
        usersRepo: usersRepo,
        sessMgr: sessMgr,
    }
}

func (rs *UsersService) RegisterUser(login string, password string) error {
    return rs.usersRepo.PostUser(uuid.New(), login, password)
}

func (rs *UsersService) LoginUser(login string, password string) (string, error) {
    if user, err := rs.usersRepo.GetUserByCred(login, password); err != nil {
        return "", err
    } else if sess, err := rs.sessMgr.StartSession(user.Id); err != nil {
        return "", err
    } else {
        return sess.SessionId, nil
    }
}
