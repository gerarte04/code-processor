package sessions

import (
	"context"
	"fmt"
	"http_server/config"
	"http_server/pkg/generator"
	"http_server/repository/models"
	"http_server/usecases"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionStorage struct {
    cli *redis.Client
    ctx context.Context
    serviceCfg config.ServiceConfig
    redisCfg config.RedisConfig
}

func NewSessionStorage(serviceCfg config.ServiceConfig, redisCfg config.RedisConfig) (*SessionStorage, error) {
    cli := redis.NewClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%s", redisCfg.Host, redisCfg.Port),
        Username: redisCfg.User,
        Password: redisCfg.UserPassword,
        DB: redisCfg.DBNumber,
    })

    err := cli.Ping(context.Background()).Err()

    if err != nil {
        return nil, err
    }

    return &SessionStorage{
        cli: cli,
        ctx: context.Background(),
        serviceCfg: serviceCfg,
        redisCfg: redisCfg,
    }, nil
}

func (sm *SessionStorage) CreateSession(userId uuid.UUID) (*models.Session, error) {
    newId := generator.NewSessionId()
    err := sm.cli.Set(sm.ctx, newId, userId.String(), sm.serviceCfg.SessionLivingTime).Err()
    
    if err != nil {
        log.Printf("starting session: putting to redis: %s", err.Error())
        return nil, usecases.ErrorInternalQueryError
    }

    return &models.Session{
        UserId: userId,
        SessionId: newId,
    }, nil
}

func (sm *SessionStorage) DeleteSession(sessionId string) error {
    err := sm.cli.Del(sm.ctx, sessionId).Err()

    if err != nil {
        log.Printf("stopping session: delete from redis: %s", err.Error())
        return usecases.ErrorInternalQueryError
    }

    return nil
}

func (sm *SessionStorage) GetSession(sessionId string) (*models.Session, error) {
    res, err := sm.cli.Get(sm.ctx, sessionId).Result()

    if err == redis.Nil {
        log.Printf("getting session: redis: %s", err.Error())
        return nil, usecases.ErrorNoSessionExists
    } else {
        uuid, err := uuid.Parse(res)

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
