package service

import (
	"context"
	"github.com/golang/glog"
	"jianghai-hu/wallet-service/internal/common"
	"jianghai-hu/wallet-service/internal/config"
	"jianghai-hu/wallet-service/utils"
)

type IMoneyActionService interface {
	Process(ctx context.Context) error
}

func GetMoneyActionService(ctx context.Context, actionType, userId, amount, orderType int) (IMoneyActionService, error) {
	switch config.DefaultMoneyActionServiceInitOption {
	case common.MONEY_ACTION_SERVICE_INIT_OPTION_DB_TRANSACTION:
		return newDBTransacMoneyActionService(userId, amount, actionType, orderType), nil
	case common.MONEY_ACTION_SERVICE_INIT_OPTION_TCC:
		return newTccMoneyActionService(userId, amount, actionType, orderType), nil
	case common.MONEY_ACTION_SERVICE_INIT_OPTION_ASYNC:
		return newAsyncTransacMoneyActionService(userId, amount, actionType, orderType), nil
	default:
		glog.ErrorContextf(ctx, "invalid MONEY_ACTION_SERVICE_INIT_OPTION:%d",
			config.DefaultMoneyActionServiceInitOption)
		return nil, utils.NewMyError(common.Constant_ERROR_SERVICE_INTERNAL, "invalid MONEY_ACTION_SERVICE_INIT_OPTION")
	}
}
