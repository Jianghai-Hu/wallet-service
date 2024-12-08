package service

import (
	"context"
	"github.com/golang/glog"
	"jianghai-hu/wallet-service/internal/common"
	"jianghai-hu/wallet-service/internal/dao"
	"jianghai-hu/wallet-service/utils"
)

type tccMoneyActionService struct {
	userId         int
	amount         int
	actionType     int
	orderType      int
	oppoUserId     int
	transactionId  int64 // filled after try phase
	walletDao      dao.IWalletDao
	transactionDao dao.ITransactionDao
}

var _ IMoneyActionService = (*tccMoneyActionService)(nil)

func newTccMoneyActionService(userId, oppoUserId, amount, actionType, orderType int) *tccMoneyActionService {
	return &tccMoneyActionService{
		userId:         userId,
		amount:         amount,
		actionType:     actionType,
		orderType:      orderType,
		oppoUserId:     oppoUserId,
		walletDao:      dao.GetWalletDao(),
		transactionDao: dao.GetTransactionDao(),
	}
}

// the implementation of tcc version is incomplete, still need to handle issues like:
// retry, empty rollback, idempotency etc.
func (service *tccMoneyActionService) Process(ctx context.Context) error {
	err := service.try(ctx)
	if err != nil {
		glog.ErrorContextf(ctx, "tccMoneyActionService|Process try phase failed:%v", err)
		_ = service.cancel(ctx) // TODO: need retry on error here
		return err
	}

	_ = service.confirm(ctx) // TODO: need retry on error here
	return nil
}

func (service *tccMoneyActionService) try(ctx context.Context) error {
	transactionId, err := utils.GetIDGenerator(ctx).Generate()
	if err != nil {
		glog.ErrorContextf(ctx, "tccMoneyActionService|gen transaction id failed, err:%v", err)
		return err
	}

	err = service.transactionDao.CreateTransaction(ctx, int64(transactionId),
		service.orderType, service.actionType, service.amount, service.userId, service.oppoUserId)
	if err != nil {
		glog.ErrorContextf(ctx, "tccMoneyActionService|CreateTransaction failed, err:%v", err)
		return err
	}
	service.transactionId = int64(transactionId)

	err = service.walletDao.FreezeBalance(ctx, service.userId, service.amount, service.actionType)
	if err != nil {
		glog.ErrorContextf(ctx, "tccMoneyActionService|FreezeBalance failed, err:%v", err)
		return err
	}
	return nil
}

func (service *tccMoneyActionService) confirm(ctx context.Context) error {
	err := service.walletDao.ConfirmBalance(ctx, service.userId, service.amount, service.actionType)
	if err != nil {
		glog.ErrorContextf(ctx, "tccMoneyActionService|ConfirmBalance failed, err:%v", err)
		return err
	}

	err = service.transactionDao.UpdateTransactionStatus(ctx, service.transactionId, common.TRANSACTION_STATUS_COMPLETE)
	if err != nil {
		glog.ErrorContextf(ctx, "tccMoneyActionService|UpdateTransactionStatus failed, err:%v", err)
		return err
	}
	return nil
}

func (service *tccMoneyActionService) cancel(ctx context.Context) error {
	if service.transactionId == 0 {
		return nil
	}

	// TODO: make this step idempotent
	err := service.walletDao.RollbackBalance(ctx, service.userId, service.amount, service.actionType)
	if err != nil {
		glog.ErrorContextf(ctx, "tccMoneyActionService|RollbackBalance failed, err:%v", err)
		return err
	}

	err = service.transactionDao.UpdateTransactionStatus(ctx, service.transactionId, common.TRANSACTION_STATUS_FAILED)
	if err != nil {
		glog.ErrorContextf(ctx, "tccMoneyActionService|UpdateTransactionStatus failed, err:%v", err)
		return err
	}
	return nil
}
