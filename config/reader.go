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
    SessionLivingTime time.Duration `yaml:"session_living_time" env-default:"30m"`
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
