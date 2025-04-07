package config

import (
	"cpapp/pkg/config/types"
	"time"
)

type ProcessorConfig struct {
    ContainerName string `yaml:"container_name" env-default:"code_container"`
    ImageName string `yaml:"image_name" env-default:"processing_code_image"`
    CodeFileName string `yaml:"code_file_name" env-default:"file"`
    ImagePath string `yaml:"image_path" env-default:"./build"`
    Dockerfile string `yaml:"dockerfile" env-default:"Dockerfile"`
    BuildTimeout time.Duration `yaml:"build_timeout" env-default:"5m"`
    RunTimeout time.Duration `yaml:"run_timeout" env-default:"10m"`
}

type Config struct {
    RabbMQCfg types.RabbitMQConfig `yaml:"rabbitmq"`
    PostgresCfg types.PostgreSQLConfig `yaml:"postgres"`
    ProcCfg ProcessorConfig `yaml:"processor"`
    PromCfg types.PrometheusConfig `yaml:"prometheus"`
}
