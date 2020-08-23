package handlers

import "github.com/saskamegaprogrammist/userBalanceService/useCases"

type Handlers struct {
	FundsHandlers    *FundsHandlers
}

var h Handlers

func Init(fundsUC useCases.FundsUCInterface) error {
	h.FundsHandlers = &FundsHandlers{fundsUC}
	return nil
}

func GetUFundsH() *FundsHandlers {
	return h.FundsHandlers
}

