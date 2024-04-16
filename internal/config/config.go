package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		Name     string `yaml:"name" env:"DB_NAME"`
	} `yaml:"database"`
	Server struct {
		Address string `yaml:"address" env:"SERVER_ADDRESS"`
		Port    string `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Redis struct {
		Host string `yaml:"host" env:"DB_HOST"`
		Port int    `yaml:"port" env:"DB_PORT"`
	} `yaml:"redis"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
