package service

import (
	"http_server/pkg/process"
	"http_server/repository"
	"http_server/repository/models"
	"http_server/usecases"
	"time"

	"github.com/google/uuid"
)

type Object struct {
    repo repository.Object
    sessMgr usecases.SessionManager
}

func NewObject(repo repository.Object, sessMgr usecases.SessionManager) *Object {
    return &Object{
        repo: repo,
        sessMgr: sessMgr,
    }
}

func (rs *Object) GetTask(key uuid.UUID, sessionId string) (*models.Task, error) {
    _, err := rs.sessMgr.GetSession(sessionId)

    if err != nil {
        return nil, err
    }

    task, err := rs.repo.GetTask(key)

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
    err = rs.repo.PostTask(key)

    if err != nil {
        return nil, err
    }

    tsk, _ := rs.repo.GetTask(key)
    go process.SleepAndComplete(tsk, dur)

    return &key, nil
}

func (rs *Object) RegisterUser(login string, password string) error {
    return rs.repo.PostUser(uuid.New(), login, password)
}

func (rs *Object) LoginUser(login string, password string) (string, error) {
    if user, err := rs.repo.GetUserByCred(login, password); err != nil {
        return "", err
    } else if sess, err := rs.sessMgr.StartSession(user.Id, time.Now().Add(time.Minute)); err != nil {
        return "", err
    } else {
        return sess.SessionId, nil
    }
}
