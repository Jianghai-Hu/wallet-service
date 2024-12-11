package config

import (
	"github.com/go-redis/redis/v8"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

var DefaultDBConfig = DBConfig{
	Host:     "localhost",
	Port:     "5432",
	Username: "jianghai",
	Password: "123456",
	DBName:   "demo_db",
}

var DefaultRedisConfig = &redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}
