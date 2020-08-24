package useCases

import "github.com/saskamegaprogrammist/userBalanceService/models"

type FundsUCInterface interface {
	Add(tx *models.Transaction) (bool, error)
	Withdraw(tx *models.Transaction) (bool, bool, error)
	Get(balance *models.Balance) (bool, error)
	Transfer(tx *models.Transaction) (bool, bool, error)
	GetTransactions(user *models.UserId, limit int, since string, sort string, desc bool) (bool, []models.Transaction, error)
}
