package models

import (
	"time"
)

//easyjson:json
type Transaction struct {
	Id            int       `json:"-"`
	UserId        int       `json:"user_id"`
	UserFromId    int       `json:"user_from_id"`
	OperationType int       `json:"operation_type"`
	Sum           float64   `json:"sum"`
	Balance       float64   `json:"-"`
	BalanceFrom   float64   `json:"-"`
	Created       time.Time `json:"created"`
}

//easyjson:json
type Transactions []Transaction
