package service

import (
	"http_server/repository"
	"http_server/repository/models"
	"http_server/usecases"
	"http_server/usecases/process"
	"time"

	"github.com/google/uuid"
)

type Object struct {
	tasksRepo repository.TasksRepo
	usersRepo repository.UsersRepo
    sessMgr usecases.SessionManager
}

func NewObject(tasksRepo repository.TasksRepo, usersRepo repository.UsersRepo, sessMgr usecases.SessionManager) *Object {
    return &Object{
        tasksRepo: tasksRepo,
        usersRepo: usersRepo,
        sessMgr: sessMgr,
    }
}

func (rs *Object) GetTask(key uuid.UUID, sessionId string) (*models.Task, error) {
    _, err := rs.sessMgr.GetSession(sessionId)

    if err != nil {
        return nil, err
    }

    task, err := rs.tasksRepo.GetTask(key)

    if err != nil {
        return nil, err
    }

    return task, nil
}

func (rs *Object) PostTask(dur time.Duration, sessionId string) (*uuid.UUID, error) {
    _, err := rs.sessMgr.GetSession(sessionId)

    if err != nil {
        return nil, err
    }

    key := uuid.New()
    err = rs.tasksRepo.PostTask(key)

    if err != nil {
        return nil, err
    }

    tsk, _ := rs.tasksRepo.GetTask(key)
    go process.SleepAndComplete(tsk, dur)

    return &key, nil
}

func (rs *Object) RegisterUser(login string, password string) error {
    return rs.usersRepo.PostUser(uuid.New(), login, password)
}

func (rs *Object) LoginUser(login string, password string) (string, error) {
    if user, err := rs.usersRepo.GetUserByCred(login, password); err != nil {
        return "", err
    } else if sess, err := rs.sessMgr.StartSession(user.Id, time.Now().Add(30 * time.Minute)); err != nil {
        return "", err
    } else {
        return sess.SessionId, nil
    }
}
