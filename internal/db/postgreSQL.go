package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang/glog"
	_ "github.com/lib/pq" // required by database/sql

	"jianghai-hu/wallet-service/internal/config"
)

var globalDB *sql.DB

func InitPostgres(ctx context.Context) {
	var err error

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DefaultDBConfig.Host,
		config.DefaultDBConfig.Port,
		config.DefaultDBConfig.Username,
		config.DefaultDBConfig.Password,
		config.DefaultDBConfig.DBName,
	)

	globalDB, err = sql.Open("postgres", connStr)
	if err != nil {
		glog.FatalContextf(ctx, "Failed to connect to Postgres: %v", err)
	}

	if err = globalDB.Ping(); err != nil {
		glog.FatalContextf(ctx, "Failed to ping Postgres: %v", err)
	}

	glog.InfoContext(ctx, "Connected to Postgres")
}

func GetDBClient(ctx context.Context) *sql.DB {
	if globalDB == nil {
		glog.FatalContext(ctx, "init DB client first!")
	}

	return globalDB
}
