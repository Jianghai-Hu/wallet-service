package config

import (
	"github.com/go-redis/redis/v8"
	"jianghai-hu/wallet-service/internal/common"
)

type DbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

var DefaultDBConfig = DbConfig{
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

var DefaultMoneyActionServiceInitOption = common.MONEY_ACTION_SERVICE_INIT_OPTION_DB_TRANSACTION
