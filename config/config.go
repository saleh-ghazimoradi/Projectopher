package config

import (
	"github.com/caarlos0/env/v11"
	"sync"
	"time"
)

var (
	instance *Config
	once     sync.Once
	initErr  error
)

type Config struct {
	Server  Server
	MongoDB MongoDB
}

type Server struct {
	Host         string        `env:"SERVER_HOST"`
	Port         string        `env:"SERVER_PORT"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`
}

type MongoDB struct {
	Host        string        `env:"MONGODB_HOST"`
	Port        string        `env:"MONGODB_PORT"`
	User        string        `env:"MONGODB_USER"`
	Pass        string        `env:"MONGODB_PASS"`
	DBName      string        `env:"MONGODB_DBNAME"`
	AuthSource  string        `env:"MONGODB_AUTH_SOURCE"`
	MaxPoolSize uint64        `env:"MONGODB_MAX_POOL_SIZE"`
	MinPoolSize uint64        `env:"MONGODB_MIN_POOL_SIZE"`
	Timeout     time.Duration `env:"MONGODB_TIMEOUT"`
}

func GetInstance() (*Config, error) {
	once.Do(func() {
		instance = &Config{}
		initErr = env.Parse(instance)
		if initErr != nil {
			initErr = nil
		}
	})
	return instance, initErr
}
