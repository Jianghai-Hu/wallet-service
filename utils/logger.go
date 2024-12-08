package utils

import (
	"context"
	"flag"
	"github.com/golang/glog"
)

func InitLogger(ctx context.Context) {
	flag.Parse()
	defer glog.Flush()
	glog.InfoContext(ctx, "InitLogger success")
	return
}
