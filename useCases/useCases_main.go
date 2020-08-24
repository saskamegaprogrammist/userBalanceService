package useCases

import "github.com/saskamegaprogrammist/userBalanceService/repository"

type UseCases struct {
	FundsUC *FundsUC
}

var uc UseCases

func Init(balanceRepo repository.BalanceRepoI, transactionsRepo repository.TransactionsRepoI) error {
	uc.FundsUC = &FundsUC{balanceRepo, transactionsRepo}
	return nil
}

func GetFundsUC() FundsUCInterface {
	return uc.FundsUC
}
