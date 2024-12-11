package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"jianghai-hu/wallet-service/internal/processor"

	"jianghai-hu/wallet-service/internal/common"
	"jianghai-hu/wallet-service/internal/db"
	"jianghai-hu/wallet-service/utils"
)

func main() {
	ctx := context.Background()
	setUp(ctx)

	router := registerAPI()
	run(ctx, router)
}

func setUp(ctx context.Context) {
	utils.InitLogger(ctx)
	db.InitPostgres(ctx)
	// cache.InitRedis(ctx)
	utils.InitIDGenerator(ctx, common.ID_GENERATOR_MACHINE_ID)
}

func registerAPI() *mux.Router {
	router := mux.NewRouter()
	for _, p := range processor.AllProcessorConfigs() {
		router.HandleFunc(p.Command, p.Processor).Methods(strings.ToUpper(p.Method))
	}

	return router
}

//nolint:mnd // fix in future
func run(ctx context.Context, router *mux.Router) {
	server := &http.Server{
		Addr:              "localhost:8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	glog.InfoContextf(ctx, "starting server on %s", server.Addr)
	glog.FatalContext(ctx, server.ListenAndServe())
}
