package useCases

import (
	"fmt"
	"github.com/saskamegaprogrammist/userBalanceService/models"
	"github.com/saskamegaprogrammist/userBalanceService/repository"
	"github.com/saskamegaprogrammist/userBalanceService/utils"
	"time"
)

type FundsUC struct {
	BalanceRepo repository.BalanceRepoI
	TransactionsRepo repository.TransactionsRepoI
}

func (fundsUC *FundsUC) Add(tx *models.Transaction) (bool, error) {
	var newBalance models.Balance
	newBalance.UserId = tx.UserId
	errType, err := fundsUC.BalanceRepo.GetBalanceByUserId(&newBalance)
	if err != nil {
		if errType == utils.USER_ERROR {
			err := fundsUC.BalanceRepo.InsertUser(&newBalance)
			if err != nil {
				return false, err
			}
		} else if errType == utils.SERVER_ERROR {
			return false, err
		}
	}

	if tx.Sum <= 0 {
		return true, fmt.Errorf("sum must be positive")
	}

	tx.Balance = newBalance.Balance + tx.Sum
	tx.OperationType = utils.GetOperationType("Add")
	tx.Created = time.Now()

	err = fundsUC.TransactionsRepo.Add(tx)
	return false, err
}

func (fundsUC *FundsUC) Withdraw(tx *models.Transaction) (bool, bool, error) {
	var newBalance models.Balance
	newBalance.UserId = tx.UserId
	errType, err := fundsUC.BalanceRepo.GetBalanceByUserId(&newBalance)
	if err != nil {
		if errType == utils.USER_ERROR {
			err := fundsUC.BalanceRepo.InsertUser(&newBalance)
			if err != nil {
				return false, false, err
			}
		} else if errType == utils.SERVER_ERROR {
			return false, false, err
		}
	}

	if tx.Sum <= 0 {
		return true, false, fmt.Errorf("sum must be positive")
	}

	tx.Balance = newBalance.Balance - tx.Sum

	if tx.Balance < 0 {
		return false, true, fmt.Errorf("you don't have enough funds")
	}

	tx.OperationType = utils.GetOperationType("Withdraw")
	tx.Created = time.Now()

	err = fundsUC.TransactionsRepo.Add(tx)
	return false, false, err
}

func (fundsUC *FundsUC) Get(balance *models.Balance) (bool, error) {
	if balance.UserId == utils.ERROR_ID {
		return true, fmt.Errorf("incorrect user id")
	}
	errType, err := fundsUC.BalanceRepo.GetBalanceByUserId(balance)
	if err != nil {
		if errType == utils.USER_ERROR {
			err := fundsUC.BalanceRepo.InsertUser(balance)
			if err != nil {
				return false, err
			}
			return false, nil
		} else if errType == utils.SERVER_ERROR {
			return false, err
		}
	}
	return false, nil
}

func (fundsUC *FundsUC) Transfer(tx *models.Transaction) (bool, bool, error) {
	var newBalance models.Balance
	newBalance.UserId = tx.UserId
	errType, err := fundsUC.BalanceRepo.GetBalanceByUserId(&newBalance)
	if err != nil {
		if errType == utils.USER_ERROR {
			err := fundsUC.BalanceRepo.InsertUser(&newBalance)
			if err != nil {
				return false, false, err
			}
		} else if errType == utils.SERVER_ERROR {
			return false, false, err
		}
	}

	var newBalanceFrom models.Balance
	newBalanceFrom.UserId = tx.UserFromId
	errType, err = fundsUC.BalanceRepo.GetBalanceByUserId(&newBalanceFrom)
	if err != nil {
		if errType == utils.USER_ERROR {
			err := fundsUC.BalanceRepo.InsertUser(&newBalanceFrom)
			if err != nil {
				return false, false, err
			}
		} else if errType == utils.SERVER_ERROR {
			return false, false, err
		}
	}

	if tx.Sum <= 0 {
		return true, false, fmt.Errorf("sum must be positive")
	}

	tx.BalanceFrom = newBalanceFrom.Balance - tx.Sum

	if tx.BalanceFrom < 0 {
		return false, true, fmt.Errorf("user doesn't have enough funds")
	}

	tx.Balance = newBalance.Balance + tx.Sum

	tx.OperationType = utils.GetOperationType("Transfer")
	tx.Created = time.Now()

	err = fundsUC.TransactionsRepo.Add(tx)
	return false, false, err
}

func (fundsUC *FundsUC) GetTransactions(user *models.UserId, limit int, since string, sort string, desc bool) (bool, []models.Transaction, error) {
	txs := make([]models.Transaction, 0)
	if user.UserId == utils.ERROR_ID {
		return true, txs, fmt.Errorf("incorrect user id")
	}
	var newBalance models.Balance
	newBalance.UserId = user.UserId
	errType, err := fundsUC.BalanceRepo.GetBalanceByUserId(&newBalance)
	if err != nil {
		if errType == utils.USER_ERROR {
			return false, txs, nil
		} else if errType == utils.SERVER_ERROR {
			return false, txs, err
		}
	}
	txs, errType, err = fundsUC.TransactionsRepo.GetUserTransactions(user, limit, since, sort, desc)
	if err != nil {
		if errType == utils.USER_ERROR {
			return true, txs, err
		} else if errType == utils.SERVER_ERROR {
			return false, txs, err
		}
	}
	return false, txs, nil
}
