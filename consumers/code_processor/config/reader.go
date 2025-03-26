package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type RabbitMQConfig struct {
    Authority string `yaml:"authority" env:"BROKER_AUTHORITY" env-default:"guest:guest"`
    Host string `yaml:"host" env:"BROKER_HOST"`
    Port string `yaml:"port" env:"BROKER_PORT"`
    QueueName string `yaml:"queue_name"`
}

type PostgreSQLConfig struct {
    Host string `yaml:"host" env:"POSTGRES_HOST"`
    Port string `yaml:"port" env:"POSTGRES_PORT"`
    DB string `yaml:"db" env:"POSTGRES_DB"`
    User string `yaml:"user" env:"POSTGRES_USER"`
    Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
}

type ProcessorConfig struct {
    ContainerName string `yaml:"container_name" env-default:"code_container"`
    ImageName string `yaml:"image_name" env-default:"processing_code_image"`
    CodeFileName string `yaml:"code_file_name" env-default:"file"`
    ImagePath string `yaml:"image_path" env-default:"./build"`
    Dockerfile string `yaml:"dockerfile" env-default:"Dockerfile"`
    BuildTimeout time.Duration `yaml:"build_timeout" env-default:"5m"`
    RunTimeout time.Duration `yaml:"run_timeout" env-default:"10m"`
}

type PrometheusConfig struct {
    ListenPort string `yaml:"listen_port" env:"PROM_LISTEN_PORT"`
}

type Config struct {
    RabbMQCfg RabbitMQConfig `yaml:"rabbitmq"`
    PostgresCfg PostgreSQLConfig `yaml:"postgres"`
    ProcCfg ProcessorConfig `yaml:"processor"`
    PromCfg PrometheusConfig `yaml:"prometheus"`
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
