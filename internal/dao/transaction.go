package dao

import (
	"context"
	"fmt"
	"jianghai-hu/wallet-service/internal/common"
	"jianghai-hu/wallet-service/internal/db"
	"time"
)

type ITransactionDao interface {
	CreateTransaction(ctx context.Context, transactionID int64, orderType, transactionType, amount, userID, oppoUserID int) error
	UpdateTransactionStatus(ctx context.Context, transactionID int64, status int) error
}

func GetTransactionDao() ITransactionDao {
	return newTransactionDao()
}

var _ ITransactionDao = (*transactionDao)(nil)

type transactionDao struct{}

func newTransactionDao() *transactionDao {
	return &transactionDao{}
}

func (dao *transactionDao) CreateTransaction(ctx context.Context, transactionID int64, orderType, transactionType, amount, userID, oppoUserID int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		INSERT INTO transaction_tab (
			transaction_id, order_type, transaction_type, amount, status, user_id, oppo_user_id, create_time, update_time, last_process_time
		) VALUES (%d, %d, %d, %d, %d, $d, %d, %d, %d, %d);
	`
	sqlStr = fmt.Sprintf(sqlStr, transactionID, orderType, transactionType, amount, common.TRANSACTION_STATUS_PENDING, userID, oppoUserID,
		timeNow, timeNow, timeNow)
	_, err := db.DBClient(ctx).Exec(sqlStr)
	return err
}

func (dao *transactionDao) UpdateTransactionStatus(ctx context.Context, transactionID int64, status int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE transaction_tab
		SET status = %d, updated_time = %d, last_process_time = %d
		WHERE transaction_id = %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, status, timeNow, timeNow, transactionID)
	_, err := db.DBClient(ctx).Exec(sqlStr)
	return err
}
