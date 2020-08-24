package repository

import "github.com/saskamegaprogrammist/userBalanceService/models"

type TransactionsRepoI interface {
	Add(transaction *models.Transaction) error
	GetUserTransactions(user *models.UserId, limit int, since string, sort string, desc bool) ([]models.Transaction, int, error)
}
