//nolint:revive,stylecheck //fix in future
package tcc_manager

import (
	"context"

	"jianghai-hu/wallet-service/internal/common"
)

type tccStatusRecorder struct {
	tryOnce     bool
	cancelOnce  bool
	confirmOnce bool
}

func newTCCStatusRecorder() *tccStatusRecorder {
	return &tccStatusRecorder{}
}

type ITCCManager interface {
	Try(ctx context.Context) error
	Confirm(ctx context.Context) error
	Cancel(ctx context.Context) error
	ReportStatus() string
}

//nolint:ireturn // fix in future
func NewTCCMangerByOrderType(ctx context.Context, userID, oppoUserID, amount, orderType int) ITCCManager {
	switch orderType {
	case common.ORDER_TYPE_WITHDRAW:
		return newTCCMoneyActionManager(ctx, userID, oppoUserID, amount, common.MONEY_ACTION_TYPE_MONEY_OUT, orderType)
	case common.ORDER_TYPE_DEPOSIT:
		return newTCCMoneyActionManager(ctx, userID, oppoUserID, amount, common.MONEY_ACTION_TYPE_MONEY_IN, orderType)
	case common.ORDER_TYPE_TRANSFER:
		return newTCCTransferManager(ctx, userID, oppoUserID, amount)
	default:
		return nil
	}
}
