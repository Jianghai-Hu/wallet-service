package dao

import (
	"context"
	"fmt"
	"jianghai-hu/wallet-service/internal/common"
	"jianghai-hu/wallet-service/internal/db"
	"jianghai-hu/wallet-service/utils"
	"time"
)

type IWalletDao interface {
	CreateWallet(ctx context.Context, userId int32) error
	FreezeBalance(ctx context.Context, userId, amount, actionType int) error
	RollbackBalance(ctx context.Context, userId, amount, actionType int) error
	ConfirmBalance(ctx context.Context, userId, amount, actionType int) error
}

func GetWalletDao() IWalletDao {
	return newWalletDao()
}

var _ IWalletDao = (*walletDao)(nil)

type walletDao struct{}

func newWalletDao() *walletDao {
	return &walletDao{}
}

func (dao *walletDao) CreateWallet(ctx context.Context, userId int32) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		INSERT INTO wallet_tab (
			user_id, balance, frozen_balance, ext_info, create_time, update_time
		) VALUES (%d, 0, 0, %s, %d, %d);
	`
	sqlStr = fmt.Sprintf(sqlStr, userId, common.DEFAULT_WALLET_EXT_INFO, timeNow, timeNow)

	_, err := db.DBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}
	return nil
}

func (dao *walletDao) FreezeBalance(ctx context.Context, userId, amount, actionType int) error {
	if actionType == common.MONEY_ACTION_TYPE_MONEY_IN {
		return dao.freezeIncome(ctx, userId, amount)
	}
	return dao.freezeDeduct(ctx, userId, amount)
}

func (dao *walletDao) freezeIncome(ctx context.Context, userID int, amount int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE wallet_tab 
		SET frozen_balance = frozen_balance + %d, update_time = %d 
		WHERE user_id = %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, amount, timeNow, userID)

	_, err := db.DBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}
	return nil
}

func (dao *walletDao) freezeDeduct(ctx context.Context, userID int, amount int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE wallet_tab 
		SET balance = balance - %d, frozen_balance = frozen_balance + %d, update_time = %d
		WHERE user_id = %d and balance >= %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, amount, amount, timeNow, userID, amount)

	result, err := db.DBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NewMyError(common.Constant_ERROR_INSUFFIENT_BALANCE,
			"freezeDeduct failed|insufficient balance")
	}
	return nil
}

func (dao *walletDao) RollbackBalance(ctx context.Context, userID int, amount, actionType int) error {
	if actionType == common.MONEY_ACTION_TYPE_MONEY_IN {
		return dao.rollbackIncome(ctx, userID, amount)
	}
	return dao.rollbackDeduct(ctx, userID, amount)
}

func (dao *walletDao) rollbackIncome(ctx context.Context, userID int, amount int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE wallet_tab 
		SET frozen_balance = frozen_balance - %d, update_time = %d
		WHERE user_id = %d and frozen_balance >= %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, amount, timeNow, userID, amount)

	result, err := db.DBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NewMyError(common.Constant_ERROR_SERVICE_INTERNAL,
			"rollbackIncome failed|insufficient frozen balance")
	}
	return err
}

func (dao *walletDao) rollbackDeduct(ctx context.Context, userID int, amount int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE wallet_tab 
		SET balance = balance + %d, frozen_balance = frozen_balance - %d, update_time = %d
		WHERE user_id = %d and frozen_balance >= %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, amount, amount, timeNow, userID, amount)

	result, err := db.DBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NewMyError(common.Constant_ERROR_SERVICE_INTERNAL,
			"rollbackDeduct failed|insufficient frozen balance")
	}
	return err
}

func (dao *walletDao) ConfirmBalance(ctx context.Context, userID, amount, actionType int) error {
	if actionType == common.MONEY_ACTION_TYPE_MONEY_IN {
		return dao.confirmIncome(ctx, userID, amount)
	}
	return dao.confirmDeduct(ctx, userID, amount)
}

func (dao *walletDao) confirmIncome(ctx context.Context, userID, amount int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE wallet_tab
		SET balance = balance + %d, frozen_balance = frozen_balance - %d, update_time = %d
		WHERE user_id = %d and frozen_balance >= %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, amount, amount, timeNow, userID, amount)

	result, err := db.DBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NewMyError(common.Constant_ERROR_SERVICE_INTERNAL,
			"confirmIncome failed|insufficient frozen balance")
	}
	return nil
}

func (dao *walletDao) confirmDeduct(ctx context.Context, userID, amount int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE wallet_tab
		SET frozen_balance = frozen_balance - %d, update_time = %d
		WHERE user_id = %d and frozen_balance >= %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, amount, timeNow, userID, amount)

	result, err := db.DBClient(ctx).Exec(sqlStr)
	if err != nil {
		return utils.WrapMyError(common.Constant_ERROR_SERVICE_INTERNAL, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return utils.NewMyError(common.Constant_ERROR_SERVICE_INTERNAL,
			"confirmDeduct failed|insufficient frozen balance")
	}
	return nil
}
