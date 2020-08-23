package repository

import (
	"fmt"
	"github.com/google/logger"
	"github.com/jackc/pgx"
	"github.com/saskamegaprogrammist/userBalanceService/models"
	"github.com/saskamegaprogrammist/userBalanceService/utils"
	"time"
)

type TransactionsRepo struct {
}

func (transactionsRepo *TransactionsRepo) Add(tx *models.Transaction) error {
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("Failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}

	row := transaction.QueryRow(`INSERT INTO transactions (user_id, user_from_id, operation, sum, balance, balance_from, created) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) returning id`,
		tx.UserId, tx.UserFromId, tx.OperationType, tx.Sum, tx.Balance, tx.BalanceFrom, tx.Created)
	err = row.Scan(&tx.Id)
	if err != nil {
		logger.Errorf("Failed to scan row: %v", err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			logger.Errorf("Failed to rollback: %v", err)
			return errRollback
		}
		return err
	}

	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("Error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}
	return nil
}

func (transactionsRepo *TransactionsRepo) GetUserTransactions(user *models.UserId, limit int, since string, sort string, desc bool) ([]models.Transaction, int, error) {
	txs := make([]models.Transaction, 0)
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("Failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return txs, utils.SERVER_ERROR, dbError
	}
	var sinceTime time.Time
	if since != "" {
		sinceTime, err = time.Parse(time.RFC3339Nano, since)
		if err != nil {
			parseError := fmt.Errorf("Failed to parse since param: %v", err.Error())
			logger.Errorf(parseError.Error())
			return txs, utils.USER_ERROR, parseError
		}
	}
	var rows *pgx.Rows
	if desc {
		if limit != utils.LIMIT_DEFAULT {
			if since != "" {
				if sort == "date" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1 AND created <= $2
							ORDER BY created DESC LIMIT $3`, user.UserId, sinceTime, limit)
				} else if sort == "sum" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1 AND created <= $2
							ORDER BY sum DESC LIMIT $3`, user.UserId, sinceTime, limit)
				} else {
					userError := fmt.Errorf("Wrong sort param")
					logger.Errorf(userError.Error())
					return txs, utils.USER_ERROR, userError
				}
			} else {
				if sort == "date" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1
							ORDER BY created DESC LIMIT $2`, user.UserId, limit)
				} else if sort == "sum" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1
							ORDER BY sum DESC LIMIT $2`, user.UserId, limit)
				} else {
					userError := fmt.Errorf("Wrong sort param")
					logger.Errorf(userError.Error())
					return txs, utils.USER_ERROR, userError
				}
			}
		} else {
			if since != "" {
				if sort == "date" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1 AND created <= $2
							ORDER BY created DESC`, user.UserId, sinceTime)
				} else if sort == "sum" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1 AND created <= $2
							ORDER BY sum DESC`, user.UserId, sinceTime)
				} else {
					userError := fmt.Errorf("Wrong sort param")
					logger.Errorf(userError.Error())
					return txs, utils.USER_ERROR, userError
				}
			} else {
				if sort == "date" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1
							ORDER BY created DESC`, user.UserId)
				} else if sort == "sum" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1
							ORDER BY sum DESC `, user.UserId)
				} else {
					userError := fmt.Errorf("Wrong sort param")
					logger.Errorf(userError.Error())
					return txs, utils.USER_ERROR, userError
				}
			}
		}
	} else {
		if limit != utils.LIMIT_DEFAULT {
			if since != "" {
				if sort == "date" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1 AND created >= $2
							ORDER BY created LIMIT $3`, user.UserId, sinceTime, limit)
				} else if sort == "sum" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1 AND created >= $2
							ORDER BY sum LIMIT $3`, user.UserId, sinceTime, limit)
				} else if sort == "" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1 AND created >= $2
							LIMIT $3`, user.UserId, sinceTime, limit)
				} else {
					userError := fmt.Errorf("Wrong sort param")
					logger.Errorf(userError.Error())
					return txs, utils.USER_ERROR, userError
				}
			} else {
				if sort == "date" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1
							ORDER BY created LIMIT $2`, user.UserId, limit)
				} else if sort == "sum" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1
							ORDER BY sum LIMIT $2`, user.UserId, limit)
				} else if sort == "" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1
							LIMIT $2`, user.UserId, limit)
				} else {
					userError := fmt.Errorf("Wrong sort param")
					logger.Errorf(userError.Error())
					return txs, utils.USER_ERROR, userError
				}
			}
		} else {
			if since != "" {
				if sort == "date" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1 AND created >= $2
							ORDER BY created DESC`, user.UserId, since)
				} else if sort == "sum" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1 AND created >= $2
							ORDER BY sum DESC`, user.UserId, since)
				} else if sort == "" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1 AND created >= $2`, user.UserId, sinceTime)
				} else {
					userError := fmt.Errorf("Wrong sort param")
					logger.Errorf(userError.Error())
					return txs, utils.USER_ERROR, userError
				}
			} else {
				if sort == "date" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1
							ORDER BY created `, user.UserId)
				} else if sort == "sum" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1
							ORDER BY sum `, user.UserId)
				} else if sort == "" {
					rows, err = transaction.Query(`SELECT * FROM transactions WHERE user_id = $1 OR user_from_id = $1`,user.UserId)
				} else {
					userError := fmt.Errorf("Wrong sort param")
					logger.Errorf(userError.Error())
					return txs, utils.USER_ERROR, userError
				}
			}
		}
	}
	if err != nil {
		logger.Errorf("Failed to retrieve transactions: %v", err)
	}
	for rows.Next() {
		var txFound models.Transaction
		err = rows.Scan(&txFound.Id, &txFound.UserId, &txFound.UserFromId, &txFound.OperationType, &txFound.Sum, &txFound.Balance, &txFound.BalanceFrom, &txFound.Created)
		if err != nil {
			logger.Errorf("Failed to retrieve transaction: %v", err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				logger.Errorf("Failed to rollback: %v", err)
				return txs, utils.SERVER_ERROR, errRollback
			}
			return txs, utils.SERVER_ERROR, err
		}
		txs = append(txs, txFound)
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("Error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return txs, utils.SERVER_ERROR, dbError
	}
	return txs, utils.NO_ERROR, nil
}