package utils

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"

	"jianghai-hu/wallet-service/internal/common"
)

func InitLogger(ctx context.Context) {
	if _, err := os.Stat(common.LOG_PATH); os.IsNotExist(err) {
		if err := os.MkdirAll(common.LOG_PATH, os.ModePerm); err != nil {
			panic(fmt.Sprintf("Failed to create log directory: %v\n", err))
		}
	}

	_ = flag.Set("log_dir", common.LOG_PATH)
	_ = flag.Set("alsologtostderr", "false")
	flag.Parse()

	defer glog.Flush()

	glog.InfoContext(ctx, "InitLogger success")
}
