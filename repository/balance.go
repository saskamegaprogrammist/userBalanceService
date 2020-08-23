package repository

import (
	"fmt"
	"github.com/google/logger"
	"github.com/saskamegaprogrammist/userBalanceService/models"
	"github.com/saskamegaprogrammist/userBalanceService/utils"
)

type BalanceRepo struct {
}


func (balanceRepo *BalanceRepo) GetBalanceByUserId(balance *models.Balance) (int, error) {
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("Failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return utils.SERVER_ERROR, dbError
	}

	row := transaction.QueryRow("SELECT id, user_id, balance::numeric FROM balance WHERE user_id = $1", balance.UserId)
	err = row.Scan(&balance.Id, &balance.UserId, &balance.Balance)
	if err != nil {
		logger.Errorf("Failed to retrieve balance: %v", err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			logger.Errorf("Failed to rollback: %v", err)
			return utils.SERVER_ERROR, errRollback
		}
		return utils.USER_ERROR, fmt.Errorf("this user doesn't exist")
	}

	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("Error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return utils.SERVER_ERROR, dbError
	}
	return utils.NO_ERROR, nil
}

func (balanceRepo *BalanceRepo) InsertUser (balance *models.Balance) error {
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("Failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}

	row := transaction.QueryRow("INSERT INTO balance (user_id) VALUES ($1) returning id",
		balance.UserId)
	err = row.Scan(&balance.Id)
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

