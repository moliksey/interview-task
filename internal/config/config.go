package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	Storage    `yaml:"storage" env-required:"true"`
	HttpServer `yaml:"http_server" env-required:"true"`
}
type HttpServer struct {
	Address     string `yaml:"address" env-default:"localhost:8080"`
	Timeout     string `yaml:"timeout" env-default:"4s"`
	IdleTimeout string `yaml:"idle-timeout" env-default:"60s"`
}
type Storage struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	Username string `yaml:"username" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres" env:"DB_PASSWORD"`
	Dbname   string `yaml:"dbname" env-default:"postgres"`
	SSLMode  string `yaml:"ssl-mode" env-default:"disable"`
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		logrus.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logrus.Fatalf("Config file %s doesn not exist", configPath)
	}
	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		logrus.Fatalf("cannot read config %s", err)
	}
	return &config
}
