package dao

import (
	"context"
	"fmt"
	"time"

	"jianghai-hu/wallet-service/internal/common"
	"jianghai-hu/wallet-service/internal/db"
)

type ITransactionDao interface {
	CreateTransaction(ctx context.Context, transactionID int64, orderType, transactionType, amount, userID, oppoUserID int) error
	UpdateTransactionStatus(ctx context.Context, transactionID int64, status int) error
}

func GetTransactionDao() *TransactionDaoImpl {
	return &TransactionDaoImpl{}
}

var _ ITransactionDao = (*TransactionDaoImpl)(nil)

type TransactionDaoImpl struct{}

func (dao *TransactionDaoImpl) CreateTransaction(ctx context.Context, transactionID int64, orderType, transactionType, amount, userID, oppoUserID int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		INSERT INTO transaction_tab (
			transaction_id, order_type, transaction_type, amount, status, user_id, oppo_user_id, create_time, update_time, last_process_time
		) VALUES (%d, %d, %d, %d, %d, %d, %d, %d, %d, %d);
	`
	sqlStr = fmt.Sprintf(sqlStr, transactionID, orderType, transactionType, amount, common.TRANSACTION_STATUS_PENDING, userID, oppoUserID,
		timeNow, timeNow, timeNow)
	_, err := db.GetDBClient(ctx).Exec(sqlStr)

	return err
}

func (dao *TransactionDaoImpl) UpdateTransactionStatus(ctx context.Context, transactionID int64, status int) error {
	timeNow := time.Now().UnixMicro()
	sqlStr := `
		UPDATE transaction_tab
		SET status = %d, update_time = %d, last_process_time = %d
		WHERE transaction_id = %d;
	`
	sqlStr = fmt.Sprintf(sqlStr, status, timeNow, timeNow, transactionID)
	_, err := db.GetDBClient(ctx).Exec(sqlStr)

	return err
}
