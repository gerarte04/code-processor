package sessions

import (
	"context"
	"http_server/config"
	"http_server/pkg/generator"
	"http_server/repository/models"
	"http_server/usecases"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionManager struct {
    db *redis.Client
    ctx context.Context
    serviceCfg config.ServiceConfig
}

func NewSessionManager(serviceCfg config.ServiceConfig) *SessionManager {
    db := redis.NewClient(&redis.Options{
        Addr: "redis:6379",
        Username: "admin-admin",
        Password: "rds",
        DB: 0,
    })

    return &SessionManager{
        db: db,
        ctx: context.Background(),
        serviceCfg: serviceCfg,
    }
}

func (sm *SessionManager) StartSession(userId uuid.UUID) (*models.Session, error) {
    newId := generator.NewSessionId()
    res := sm.db.Set(sm.ctx, newId, userId.String(), sm.serviceCfg.SessionLivingTime)
    
    if res.Err() != nil {
        log.Printf("starting session: putting to redis: %s", res.Err().Error())
        return nil, usecases.ErrorInternalQueryError
    }

    return &models.Session{
        UserId: userId,
        SessionId: newId,
    }, nil
}

func (sm *SessionManager) StopSession(sessionId string) error {
    res := sm.db.Del(sm.ctx, sessionId)

    if res.Err() != nil {
        log.Printf("stopping session: delete from redis: %s", res.Err().Error())
        return usecases.ErrorInternalQueryError
    }

    return nil
}

func (sm *SessionManager) GetSession(sessionId string) (*models.Session, error) {
    res := sm.db.Get(sm.ctx, sessionId)

    if res.Err() == redis.Nil {
        return nil, usecases.ErrorNoSessionExists
    } else {
        uuid, err := uuid.Parse(res.Val())

        if err != nil {
            log.Printf("getting session: parsing uuid: %s", err)
            return nil, usecases.ErrorInternalQueryError
        }

        return &models.Session{
            UserId: uuid,
            SessionId: sessionId,
        }, nil
    }
}
