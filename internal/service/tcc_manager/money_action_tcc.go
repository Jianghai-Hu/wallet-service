//nolint:revive,stylecheck //fix in future
package tcc_manager

import (
	"context"
	"fmt"

	"github.com/golang/glog"

	"jianghai-hu/wallet-service/internal/common"
	"jianghai-hu/wallet-service/internal/dao"
	"jianghai-hu/wallet-service/utils"
)

type tccMoneyActionManager struct {
	// biz payload
	userID        int
	amount        int
	actionType    int
	orderType     int
	oppoUserID    int
	transactionID int64
	// tcc related
	balanceStatus     *tccStatusRecorder
	transactionStatus *tccStatusRecorder
	// dao
	walletDao      dao.IWalletDao
	transactionDao dao.ITransactionDao
}

var _ ITCCManager = (*tccMoneyActionManager)(nil)

func newTCCMoneyActionManager(ctx context.Context, userID, oppoUserID, amount, actionType, orderType int) *tccMoneyActionManager {
	transactionID, err := utils.GetIDGenerator(ctx).Generate()
	if err != nil { // TODO: provide better error handling
		glog.FatalContextf(ctx, "newTccMoneyActionManager|gen transaction id failed, err:%v", err)
	}

	return &tccMoneyActionManager{
		userID:            userID,
		amount:            amount,
		actionType:        actionType,
		orderType:         orderType,
		oppoUserID:        oppoUserID,
		transactionID:     transactionID,
		balanceStatus:     newTCCStatusRecorder(),
		transactionStatus: newTCCStatusRecorder(),
		walletDao:         dao.GetWalletDao(),
		transactionDao:    dao.GetTransactionDao(),
	}
}

func (manager *tccMoneyActionManager) Try(ctx context.Context) error {
	if manager.transactionStatus.cancelOnce || manager.balanceStatus.cancelOnce {
		return utils.NewMyError(common.Constant_ERROR_SERVICE_INTERNAL,
			"Try happens after Cancel")
	}

	if !manager.transactionStatus.tryOnce {
		err := manager.transactionDao.CreateTransaction(ctx, manager.transactionID,
			manager.orderType, manager.actionType, manager.amount, manager.userID, manager.oppoUserID)
		if err != nil {
			glog.ErrorContextf(ctx, "tccMoneyActionManager|CreateTransaction failed, err:%v", err)
			return err
		}

		manager.transactionStatus.tryOnce = true
	}

	if !manager.balanceStatus.tryOnce {
		err := manager.walletDao.FreezeBalance(ctx, manager.userID, manager.amount, manager.actionType)
		if err != nil {
			glog.ErrorContextf(ctx, "tccMoneyActionManager|FreezeBalance failed, err:%v", err)
			return err
		}

		manager.balanceStatus.tryOnce = true
	}

	return nil
}

func (manager *tccMoneyActionManager) Confirm(ctx context.Context) error {
	if !manager.balanceStatus.confirmOnce {
		err := manager.walletDao.ConfirmBalance(ctx, manager.userID, manager.amount, manager.actionType)
		if err != nil {
			glog.ErrorContextf(ctx, "tccMoneyActionManager|ConfirmBalance failed, err:%v", err)
			return err
		}

		manager.balanceStatus.confirmOnce = true
	}

	if !manager.transactionStatus.confirmOnce {
		err := manager.transactionDao.UpdateTransactionStatus(ctx, manager.transactionID, common.TRANSACTION_STATUS_COMPLETE)
		if err != nil {
			glog.ErrorContextf(ctx, "tccMoneyActionManager|UpdateTransactionStatus failed, err:%v", err)
			return err
		}

		manager.transactionStatus.confirmOnce = true
	}

	return nil
}

func (manager *tccMoneyActionManager) Cancel(ctx context.Context) error {
	if !manager.balanceStatus.tryOnce { // allow empty rollback
		manager.balanceStatus.cancelOnce = true
	} else if !manager.balanceStatus.cancelOnce {
		err := manager.walletDao.RollbackBalance(ctx, manager.userID, manager.amount, manager.actionType)
		if err != nil {
			glog.ErrorContextf(ctx, "tccMoneyActionManager|RollbackBalance failed, err:%v", err)
			return err
		}

		manager.balanceStatus.cancelOnce = true
	}

	if !manager.transactionStatus.tryOnce { // allow empty rollback
		manager.transactionStatus.cancelOnce = true
	} else if !manager.transactionStatus.cancelOnce {
		err := manager.transactionDao.UpdateTransactionStatus(ctx, manager.transactionID, common.TRANSACTION_STATUS_FAILED)
		if err != nil {
			glog.ErrorContextf(ctx, "tccMoneyActionManager|UpdateTransactionStatus failed, err:%v", err)
			return err
		}

		manager.transactionStatus.cancelOnce = true
	}

	return nil
}

func (manager *tccMoneyActionManager) ReportStatus() string {
	return fmt.Sprintf("status report|userID:%d|actionType:%d|orderType:%d|amount:%d|balanceStatus:%+v|transactionStatus:%+v",
		manager.userID, manager.actionType, manager.orderType, manager.amount, manager.balanceStatus, manager.transactionStatus)
}
