//nolint:revive,stylecheck //fix in future
package tcc_manager

import (
	"context"
	"fmt"

	"github.com/golang/glog"

	"jianghai-hu/wallet-service/internal/common"
	"jianghai-hu/wallet-service/utils"
)

type tccTransferManger struct {
	// biz payload
	fromUserID int
	toUserID   int
	amount     int
	// tcc related
	moneyOutStatus    *tccStatusRecorder
	moneyInStatus     *tccStatusRecorder
	moneyOutTCCManger ITCCManager
	moneyInTCCManger  ITCCManager
}

var _ ITCCManager = (*tccTransferManger)(nil)

func newTCCTransferManager(ctx context.Context, fromUserID, toUserID, amount int) *tccTransferManger {
	return &tccTransferManger{
		fromUserID:     fromUserID,
		toUserID:       toUserID,
		amount:         amount,
		moneyInStatus:  newTCCStatusRecorder(),
		moneyOutStatus: newTCCStatusRecorder(),
		moneyInTCCManger: newTCCMoneyActionManager(ctx, toUserID, fromUserID, amount,
			common.MONEY_ACTION_TYPE_MONEY_IN, common.ORDER_TYPE_TRANSFER),
		moneyOutTCCManger: newTCCMoneyActionManager(ctx, fromUserID, toUserID, amount,
			common.MONEY_ACTION_TYPE_MONEY_OUT, common.ORDER_TYPE_TRANSFER),
	}
}

func (manager *tccTransferManger) Try(ctx context.Context) error {
	if manager.moneyOutStatus.cancelOnce || manager.moneyInStatus.cancelOnce {
		return utils.NewMyError(common.Constant_ERROR_SERVICE_INTERNAL,
			"Try happens after Cancel")
	}

	if !manager.moneyOutStatus.tryOnce {
		manager.moneyOutStatus.tryOnce = true

		err := manager.moneyOutTCCManger.Try(ctx)
		if err != nil {
			glog.ErrorContextf(ctx, "tccTransferManger|money out try phase failed:%v", err)
			return err
		}
	}

	if !manager.moneyInStatus.tryOnce {
		manager.moneyInStatus.tryOnce = true

		err := manager.moneyInTCCManger.Try(ctx)
		if err != nil {
			glog.ErrorContextf(ctx, "tccTransferManger|money in try phase failed:%v", err)
			return err
		}
	}

	return nil
}

func (manager *tccTransferManger) Confirm(ctx context.Context) error {
	if !manager.moneyOutStatus.confirmOnce {
		err := manager.moneyOutTCCManger.Confirm(ctx)
		if err != nil {
			glog.ErrorContextf(ctx, "tccTransferManger|money out confirm phase failed:%v", err)
			return err
		}

		manager.moneyOutStatus.confirmOnce = true
	}

	if !manager.moneyInStatus.confirmOnce {
		err := manager.moneyInTCCManger.Confirm(ctx)
		if err != nil {
			glog.ErrorContextf(ctx, "tccTransferManger|money in confirm phase failed:%v", err)
			return err
		}

		manager.moneyInStatus.confirmOnce = true
	}

	return nil
}

func (manager *tccTransferManger) Cancel(ctx context.Context) error {
	if !manager.moneyOutStatus.tryOnce { // allow empty cancel
		manager.moneyOutStatus.cancelOnce = true
	} else if !manager.moneyOutStatus.cancelOnce {
		err := manager.moneyOutTCCManger.Cancel(ctx)
		if err != nil {
			glog.ErrorContextf(ctx, "tccTransferManger|money out cancel phase failed:%v", err)
			return err
		}

		manager.moneyOutStatus.cancelOnce = true
	}

	if !manager.moneyInStatus.tryOnce { // allow empty cancel
		manager.moneyInStatus.cancelOnce = true
	} else if !manager.moneyInStatus.cancelOnce {
		err := manager.moneyInTCCManger.Cancel(ctx)
		if err != nil {
			glog.ErrorContextf(ctx, "tccTransferManger|money in cancel phase failed:%v", err)
			return err
		}

		manager.moneyInStatus.cancelOnce = true
	}

	return nil
}

func (manager *tccTransferManger) ReportStatus() string {
	return fmt.Sprintf("status report|fromUserID:%d|toUserID:%d|amount:%d|MoneyInStatus:%+v|MoneyOutStatus:%+v",
		manager.fromUserID, manager.toUserID, manager.amount, manager.moneyInStatus, manager.moneyOutStatus)
}
