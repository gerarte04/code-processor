package sessions

import (
	"http_server/config"
	"http_server/pkg/generator"
	"http_server/repository/models"
	"http_server/usecases"
	"time"

	"github.com/google/uuid"
)

type SessionManager struct {
    sessions map[string]*models.Session
    serviceCfg config.ServiceConfig
}

func NewSessionManager(serviceCfg config.ServiceConfig) *SessionManager {
    return &SessionManager{
        sessions: make(map[string]*models.Session),
        serviceCfg: serviceCfg,
    }
}

func (sm *SessionManager) StartSession(userId uuid.UUID) (*models.Session, error) {
    for _, v := range sm.sessions {
        if v.UserId == userId {
            if time.Now().After(v.ExpiresAt) {
                sm.StopSession(v.SessionId)
                break
            } else {
                return nil, usecases.ErrorUserSessionExists
            }
        }
    }

    newId := generator.NewSessionId()

    sm.sessions[newId] = &models.Session{
        UserId: userId,
        SessionId: newId,
        ExpiresAt: time.Now().Add(sm.serviceCfg.SessionLivingTime),
    }
    return sm.sessions[newId], nil
}

func (sm *SessionManager) StopSession(sessionId string) error {
    if _, ok := sm.sessions[sessionId]; !ok {
        return usecases.ErrorNoUserSessionExists
    }

    delete(sm.sessions, sessionId)
    return nil
}

func (sm *SessionManager) GetSession(sessionId string) (*models.Session, error) {
    if sess, ok := sm.sessions[sessionId]; !ok {
        return nil, usecases.ErrorNoSessionExists
    } else if time.Now().After(sess.ExpiresAt) {
        sm.StopSession(sess.SessionId)
        return nil, usecases.ErrorSessionExpired
    } else {
        return sess, nil
    }
}
