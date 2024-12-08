package service

import (
	"context"
	"jianghai-hu/wallet-service/internal/common"
	"jianghai-hu/wallet-service/utils"
)

type asyncTransacMoneyActionService struct {
	userId     int
	amount     int
	actionType int
	orderType  int
}

var _ IMoneyActionService = (*asyncTransacMoneyActionService)(nil)

func newAsyncTransacMoneyActionService(userId, amount, actionType, orderType int) *asyncTransacMoneyActionService {
	return &asyncTransacMoneyActionService{
		userId:     userId,
		amount:     amount,
		actionType: actionType,
		orderType:  orderType,
	}
}

func (service *asyncTransacMoneyActionService) Process(ctx context.Context) error {
	return utils.NewMyError(common.Constant_ERROR_SERVICE_INTERNAL, "async money action not implemented")
}
