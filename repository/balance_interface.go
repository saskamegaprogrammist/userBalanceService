package repository

import "github.com/saskamegaprogrammist/userBalanceService/models"

type BalanceRepoI interface {
	GetBalanceByUserId(user *models.Balance) (int, error)
	InsertUser(balance *models.Balance) error
}
