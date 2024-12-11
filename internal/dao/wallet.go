package dao

import (
	"context"
	"fmt"
	"time"

	"jianghai-hu/wallet-service/internal/common"
	"jianghai-hu/wallet-service/internal/db"
	"jianghai-hu/wallet-service/utils"
)

type IWalletDao interface {
	CreateWallet(ctx context.Context, userID int32) error
	FreezeBalance(ctx context.Context, userID, amount, actionType int) error
	RollbackBalance(ctx context.Context, userID, amount, actionType int) error
	ConfirmBalance(ctx context.Context, userID, amount, actionType int) error
}

func GetWalletDao() *WalletDaoImpl {
	return &WalletDaoImpl{}
}

var _ IWalletDao = (*WalletDaoImpl)(nil)

type WalletDaoImpl struct{}

func (dao *WalletDaoImpl) CreateWallet(ctx context.Context, userID int32) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		INSERT INTO wallet_tab (
			user_id, balance, frozen_balance, ext_info, create_time, update_time
		) VALUES (%d, 0, 0, %s, %d, %d);
	`
	sqlStr = fmt.Sprintf(sqlStr, userID, common.DEFAULT_WALLET_EXT_INFO, timeNow, timeNow)

	_, err := db.GetDBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}

	return nil
}

func (dao *WalletDaoImpl) FreezeBalance(ctx context.Context, userID, amount, actionType int) error {
	if actionType == common.MONEY_ACTION_TYPE_MONEY_IN {
		return dao.freezeIncome(ctx, userID, amount)
	}

	return dao.freezeDeduct(ctx, userID, amount)
}

func (dao *WalletDaoImpl) freezeIncome(ctx context.Context, userID int, amount int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE wallet_tab 
		SET frozen_balance = frozen_balance + %d, update_time = %d 
		WHERE user_id = %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, amount, timeNow, userID)

	result, err := db.GetDBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NewMyError(common.Constant_ERROR_PARAM,
			fmt.Sprintf("freezeDeduct failed|user_id:%d not exist", userID))
	}

	return nil
}

func (dao *WalletDaoImpl) freezeDeduct(ctx context.Context, userID int, amount int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE wallet_tab 
		SET balance = balance - %d, frozen_balance = frozen_balance + %d, update_time = %d
		WHERE user_id = %d and balance >= %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, amount, amount, timeNow, userID, amount)

	result, err := db.GetDBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NewMyError(common.Constant_ERROR_INSUFFIENT_BALANCE,
			fmt.Sprintf("freezeDeduct failed|insufficient balance|user_id:%d", userID))
	}

	return nil
}

func (dao *WalletDaoImpl) RollbackBalance(ctx context.Context, userID int, amount, actionType int) error {
	if actionType == common.MONEY_ACTION_TYPE_MONEY_IN {
		return dao.rollbackIncome(ctx, userID, amount)
	}

	return dao.rollbackDeduct(ctx, userID, amount)
}

func (dao *WalletDaoImpl) rollbackIncome(ctx context.Context, userID int, amount int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE wallet_tab 
		SET frozen_balance = frozen_balance - %d, update_time = %d
		WHERE user_id = %d and frozen_balance >= %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, amount, timeNow, userID, amount)

	result, err := db.GetDBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NewMyError(common.Constant_ERROR_SERVICE_INTERNAL,
			fmt.Sprintf("rollbackIncome failed|insufficient frozen balance|userID:%d", userID))
	}

	return nil
}

func (dao *WalletDaoImpl) rollbackDeduct(ctx context.Context, userID int, amount int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE wallet_tab 
		SET balance = balance + %d, frozen_balance = frozen_balance - %d, update_time = %d
		WHERE user_id = %d and frozen_balance >= %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, amount, amount, timeNow, userID, amount)

	result, err := db.GetDBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NewMyError(common.Constant_ERROR_SERVICE_INTERNAL,
			fmt.Sprintf("rollbackDeduct failed|insufficient frozen balance|userID:%d", userID))
	}

	return nil
}

func (dao *WalletDaoImpl) ConfirmBalance(ctx context.Context, userID, amount, actionType int) error {
	if actionType == common.MONEY_ACTION_TYPE_MONEY_IN {
		return dao.confirmIncome(ctx, userID, amount)
	}

	return dao.confirmDeduct(ctx, userID, amount)
}

func (dao *WalletDaoImpl) confirmIncome(ctx context.Context, userID, amount int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE wallet_tab
		SET balance = balance + %d, frozen_balance = frozen_balance - %d, update_time = %d
		WHERE user_id = %d and frozen_balance >= %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, amount, amount, timeNow, userID, amount)

	result, err := db.GetDBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NewMyError(common.Constant_ERROR_SERVICE_INTERNAL,
			fmt.Sprintf("confirmIncome failed|insufficient frozen balance|userID:%d", userID))
	}

	return nil
}

func (dao *WalletDaoImpl) confirmDeduct(ctx context.Context, userID, amount int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE wallet_tab
		SET frozen_balance = frozen_balance - %d, update_time = %d
		WHERE user_id = %d and frozen_balance >= %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, amount, timeNow, userID, amount)

	result, err := db.GetDBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NewMyError(common.Constant_ERROR_SERVICE_INTERNAL,
			fmt.Sprintf("confirmDeduct failed|insufficient frozen balance|userID:%d", userID))
	}

	return nil
}
