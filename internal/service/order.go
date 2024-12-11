package service

import (
	"context"

	"github.com/golang/glog"

	"jianghai-hu/wallet-service/internal/common"
	"jianghai-hu/wallet-service/internal/service/tcc_manager"
	"jianghai-hu/wallet-service/utils"
)

type IOrderService interface {
	Deposit(ctx context.Context, userID int, amount int) error
	Withdraw(ctx context.Context, userID int, amount int) error
	Transfer(ctx context.Context, fromUserID int, toUserID int, amount int) error
}

type OrderServiceImpl struct{}

func NewOrderService() *OrderServiceImpl {
	return &OrderServiceImpl{}
}

func (service *OrderServiceImpl) Deposit(ctx context.Context, userID, amount int) error {
	if userID <= 0 || amount <= 0 {
		return utils.NewMyError(common.Constant_ERROR_PARAM, "invalid params")
	}

	manager := tcc_manager.NewTCCMangerByOrderType(ctx, userID, 0, amount, common.ORDER_TYPE_DEPOSIT)
	defer func() {
		glog.InfoContext(ctx, manager.ReportStatus())
	}()

	err := manager.Try(ctx)
	if err != nil {
		glog.ErrorContextf(ctx, "Deposit|process Try phase failed:%v", err)
		_ = manager.Cancel(ctx) // TODO: need retry on error here

		return err
	}

	_ = manager.Confirm(ctx) // TODO: need retry on error here

	return nil
}

func (service *OrderServiceImpl) Withdraw(ctx context.Context, userID, amount int) error {
	if userID <= 0 || amount <= 0 {
		return utils.NewMyError(common.Constant_ERROR_PARAM, "invalid params")
	}

	manager := tcc_manager.NewTCCMangerByOrderType(ctx, userID, 0, amount, common.ORDER_TYPE_WITHDRAW)
	defer func() {
		glog.InfoContext(ctx, manager.ReportStatus())
	}()

	err := manager.Try(ctx)
	if err != nil {
		glog.ErrorContextf(ctx, "Withdraw|process Try phase failed:%v", err)
		_ = manager.Cancel(ctx) // TODO: need retry on error here

		return err
	}

	_ = manager.Confirm(ctx) // TODO: need retry on error here

	return nil
}

func (service *OrderServiceImpl) Transfer(ctx context.Context, fromUserID int, toUserID int, amount int) error {
	if fromUserID <= 0 || toUserID <= 0 || amount <= 0 {
		return utils.NewMyError(common.Constant_ERROR_PARAM, "invalid params")
	}

	manager := tcc_manager.NewTCCMangerByOrderType(ctx, fromUserID, toUserID, amount, common.ORDER_TYPE_TRANSFER)
	defer func() {
		glog.InfoContext(ctx, manager.ReportStatus())
	}()

	err := manager.Try(ctx)
	if err != nil {
		glog.ErrorContextf(ctx, "Transfer|process Try phase failed:%v", err)
		_ = manager.Cancel(ctx) // TODO: need retry on error here

		return err
	}

	_ = manager.Confirm(ctx) // TODO: need retry on error here

	return nil
}
