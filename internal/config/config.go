package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Level      string     `env:"LEVEL" yaml:"level" env-default:"local"`
	Storage    Storage    `env:"" yaml:"storage" env-required:"true"`
	Nats       Nats       `env:"" yaml:"nats" env-required:"true"`
	HTTPServer HTTPServer `env:"" yaml:"http_server"`
}

type Nats struct {
	ClusterID    string `env:"CLUSTER_ID" yaml:"cluster_id" env-required:"true"`
	ClientId     string `env:"CLIENT_ID" yaml:"client_id" env-required:"true"`
	ListenChanel string `env:"LISTEN_CHANNEL" yaml:"listen_chanel" env-required:"true"`
	ListenUrl    string `env:"LISTEN_URL" yaml:"listen_url" env-required:"true"`
}

type HTTPServer struct {
	Host string `env:"HTTP_HOST" yaml:"host" env-default:"localhost"`
	Port string `env:"HTTP_PORT" yaml:"port" env-default:"8080"`
}

type Storage struct {
	Username string `env:"DBUSER" yaml:"username" env-default:"postgres"`
	Password string `env:"DBPASSWORD" yaml:"password" env-required:"true"`
	Database string `env:"DATABASE" yaml:"database" env-required:"true"`
	Host     string `env:"DBHOST" yaml:"host" env-default:"localhost"`
	Port     string `env:"DBPORT" yaml:"port" env-default:"5432"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		return nil
	}
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil
	}

	return cfg
}
