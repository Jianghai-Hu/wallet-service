package server

import (
	"context"
	"jianghai-hu/wallet-service/internal/cache"
	"jianghai-hu/wallet-service/internal/common"
	"jianghai-hu/wallet-service/internal/db"
	"jianghai-hu/wallet-service/utils"
)

func main() {
	ctx := context.Background()
	utils.InitLogger(ctx)

	db.InitPostgres(ctx)
	cache.InitRedis(ctx)
	utils.InitIDGenerator(ctx, common.ID_GENERATOR_MACHINE_ID)
}
