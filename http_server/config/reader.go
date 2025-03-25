package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpConfig struct {
    Host string `yaml:"host" env:"HTTP_HOST"`
    Port string `yaml:"port" env:"HTTP_PORT"`
}

type ServiceConfig struct {
    SessionLivingTime time.Duration `yaml:"session_living_time" env-default:"30m"`
}

type PostgreSQLConfig struct {
    Host string `yaml:"host" env:"POSTGRES_HOST"`
    Port string `yaml:"port" env:"POSTGRES_PORT"`
    DB string `yaml:"db" env:"POSTGRES_DB"`
    User string `yaml:"user" env:"POSTGRES_USER"`
    Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
}

type RedisConfig struct {
    Host string `yaml:"host" env:"REDIS_HOST"`
    Port string `yaml:"port" env:"REDIS_PORT"`
    Password string `yaml:"password" env:"REDIS_PASSWORD"`
    User string `yaml:"user" env:"REDIS_USER"`
    UserPassword string `yaml:"user_password" env:"REDIS_USER_PASSWORD"`
    DBNumber int `yaml:"db_number" env:"REDIS_DB_NUMBER" env-default:"0"`
}

type RabbitMQConfig struct {
    Authority string `yaml:"authority" env:"BROKER_AUTHORITY" env-default:"guest:guest"`
    Host string `yaml:"host" env:"BROKER_HOST"`
    Port string `yaml:"port" env:"BROKER_PORT"`
    QueueName string `yaml:"queue_name"`
}

type Config struct {
    HttpCfg HttpConfig `yaml:"http"`
    PostgresCfg PostgreSQLConfig `yaml:"postgres"`
    RedisCfg RedisConfig `yaml:"redis"`
    RabbMQCfg RabbitMQConfig `yaml:"rabbitmq"`
    ServiceCfg ServiceConfig `yaml:"service"`
}

func LoadConfig(path string, cfg any) {
    if path == "" {
        log.Fatalf("path is not set")
    } else if _, err := os.Stat(path); os.IsNotExist(err) {
        log.Fatalf("config file does not exist by this path: %s", path)
    } else if err := cleanenv.ReadConfig(path, cfg); err != nil {
        log.Fatalf("error reading config: %s", err)
    }
}
