package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Log      `yaml:"logger"`
		Database `yaml:"database"`
		Services `yaml:"services"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port      string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		JWTSecret string `env-required:"true" yaml:"jwtSecret" env:"JWT_SECRET"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
		Type  string `env-required:"true" yaml:"type"  env:"LOG_TYPE"`
	}

	Database struct {
		Driver string `env-required:"true"  yaml:"driver" env:"DATABASE_DRIVER"`
		URI    string `                                   env:"DATABASE_URI"`
	}

	Services struct {
		Logo LogoService `yaml:"logo"`
	}

	LogoService struct {
		BaseUrl string `env-required:"true" yaml:"base_url" env:"BASE_URL"`
		ApiKey  string `env-required:"true" yaml:"api_key"  env:"API_KEY"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
