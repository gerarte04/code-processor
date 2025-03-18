package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpConfig struct {
    Host string `yaml:"host"`
    Port string `yaml:"port"`
}

type ServiceConfig struct {
    SessionLivingTime time.Duration `yaml:"session_living_time" env-default:"30m"`
}

type RabbitMQConfig struct {
    Authority string `yaml:"authority" env-default:"guest:guest"`
    Host string `yaml:"host" env:"BROKER_HOST"`
    Port string `yaml:"port" env:"BROKER_PORT"`
    QueueName string `yaml:"queue_name"`
}

type Config struct {
    HttpCfg HttpConfig `yaml:"http"`
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
