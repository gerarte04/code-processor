package config

import (
	"cpapp/pkg/config/types"
	"time"
)

type ServiceConfig struct {
    SessionLivingTime time.Duration `yaml:"session_living_time" env-default:"30m"`
}

type Config struct {
    HttpCfg types.HttpConfig `yaml:"http"`
    PostgresCfg types.PostgreSQLConfig `yaml:"postgres"`
    RedisCfg types.RedisConfig `yaml:"redis"`
    RabbMQCfg types.RabbitMQConfig `yaml:"rabbitmq"`
    ServiceCfg ServiceConfig `yaml:"service"`
}
