package service

import "context"

type dbTransacMoneyActionService struct {
	userId     int
	amount     int
	actionType int
	orderType  int
}

var _ IMoneyActionService = (*dbTransacMoneyActionService)(nil)

func newDBTransacMoneyActionService(userId, amount, actionType, orderType int) *dbTransacMoneyActionService {
	return &dbTransacMoneyActionService{
		userId:     userId,
		amount:     amount,
		actionType: actionType,
		orderType:  orderType,
	}
}

func (service *dbTransacMoneyActionService) Process(ctx context.Context) error {
	return nil
}
