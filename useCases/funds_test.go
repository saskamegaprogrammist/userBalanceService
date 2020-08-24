package useCases

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/saskamegaprogrammist/userBalanceService/models"
	"github.com/saskamegaprogrammist/userBalanceService/repository"
	"github.com/saskamegaprogrammist/userBalanceService/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testTxOne = models.Transaction{
	UserId: 1,
	Sum:    100,
}

var testTxTwo = models.Transaction{
	UserId: 1,
	Sum:    50,
}

var testTxWrong = models.Transaction{
	UserId: 0,
	Sum:    1001,
}

var testTxWrongSum = models.Transaction{
	UserId: 1,
	Sum:    -1001,
}

var testUserOne = models.UserId{
	UserId: 1,
}

var testUserTwo = models.UserId{
	UserId: 2,
}

var testUserWrong = models.UserId{
	UserId: 0,
}

var testBalanceOneGet = models.Balance{
	UserId: 1,
}

var testBalanceTwoGet = models.Balance{
	UserId: 2,
}

var testBalanceWrongGet = models.Balance{
	UserId:   0,
	Currency: utils.CURRENCY,
}

var testTxOneTransfer = models.Transaction{
	UserId:     2,
	UserFromId: 1,
	Sum:        100,
}

var testTxWrongTransfer = models.Transaction{
	UserId:     2,
	UserFromId: 1,
	Sum:        -100,
}

var testTransactions = []models.Transaction{
	{Id: 1, UserId: 1, UserFromId: 2, Sum: 10, Balance: 101, BalanceFrom: 100, Created: time.Now()},
	{Id: 2, UserId: 2, UserFromId: 1, Sum: 10, Balance: 101, BalanceFrom: 100, Created: time.Now()},
}

var limitInt = utils.LIMIT_DEFAULT
var descBool = false
var since = ""
var sort = ""

func TestAddFunds(t *testing.T) {
	t.Run("FundsAddOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var testBalanceOneGetLocal = models.Balance{
			UserId: 1,
		}

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGetLocal).DoAndReturn(func(user *models.Balance) (int, error) {
			user.Balance = 1000
			return utils.NO_ERROR, nil
		})

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)
		mockRepoTxs.EXPECT().Add(&testTxOne).Return(nil)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, err := fundsUseCase.Add(&testTxOne)

		assert.NoError(t, err)
		assert.Equal(t, utils.GetOperationType("Add"), testTxOne.OperationType)
		assert.Equal(t, testTxOne.Sum+1000, testTxOne.Balance)
		assert.Equal(t, false, userError)
	})

	t.Run("FundsAddInvalidUserId", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, err := fundsUseCase.Add(&testTxWrong)

		assert.Error(t, err)
		assert.Equal(t, true, userError)
		assert.Equal(t, "incorrect user id", err.Error())

	})

	t.Run("DBErrorFirst", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.SERVER_ERROR, errors.New("db error"))

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, err := fundsUseCase.Add(&testTxOne)

		assert.Error(t, err)
		assert.Equal(t, false, userError)
	})

	t.Run("DBErrorSecond", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.USER_ERROR, errors.New("no user"))
		mockRepoBalance.EXPECT().InsertUser(&testBalanceOneGet).Return(errors.New("db error"))

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, err := fundsUseCase.Add(&testTxOne)

		assert.Error(t, err)
		assert.Equal(t, false, userError)
	})

	t.Run("WrongSum", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.USER_ERROR, errors.New("no user"))
		mockRepoBalance.EXPECT().InsertUser(&testBalanceOneGet).Return(nil)

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, err := fundsUseCase.Add(&testTxWrongSum)

		assert.Error(t, err)
		assert.Equal(t, true, userError)
		assert.Equal(t, "sum must be positive", err.Error())
	})
}

func TestWithdrawFunds(t *testing.T) {
	t.Run("FundsWithdrawOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var testBalanceOneGetLocal = models.Balance{
			UserId: 1,
		}

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGetLocal).DoAndReturn(func(user *models.Balance) (int, error) {
			user.Balance = 1000
			return utils.NO_ERROR, nil
		})

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)
		mockRepoTxs.EXPECT().Add(&testTxOne).Return(nil)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Withdraw(&testTxOne)

		assert.NoError(t, err)
		assert.Equal(t, utils.GetOperationType("Withdraw"), testTxOne.OperationType)
		assert.Equal(t, 1000-testTxOne.Sum, testTxOne.Balance)
		assert.Equal(t, false, userError)
		assert.Equal(t, false, lowFunds)
	})

	t.Run("FundsAddInvalidUserId", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Withdraw(&testTxWrong)

		assert.Error(t, err)
		assert.Equal(t, true, userError)
		assert.Equal(t, false, lowFunds)
		assert.Equal(t, "incorrect user id", err.Error())

	})

	t.Run("DBErrorFirst", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.SERVER_ERROR, errors.New("db error"))

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Withdraw(&testTxOne)

		assert.Error(t, err)
		assert.Equal(t, false, lowFunds)
		assert.Equal(t, false, userError)
	})

	t.Run("DBErrorSecond", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.USER_ERROR, errors.New("no user"))
		mockRepoBalance.EXPECT().InsertUser(&testBalanceOneGet).Return(errors.New("db error"))

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Withdraw(&testTxOne)

		assert.Error(t, err)
		assert.Equal(t, false, lowFunds)
		assert.Equal(t, false, userError)
	})

	t.Run("WrongSum", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.USER_ERROR, errors.New("no user"))
		mockRepoBalance.EXPECT().InsertUser(&testBalanceOneGet).Return(nil)

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Withdraw(&testTxWrongSum)

		assert.Error(t, err)
		assert.Equal(t, true, userError)
		assert.Equal(t, false, lowFunds)
		assert.Equal(t, "sum must be positive", err.Error())
	})

	t.Run("LowFunds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).DoAndReturn(func(user *models.Balance) (int, error) {
			user.Balance = 50
			return utils.NO_ERROR, nil
		})

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Withdraw(&testTxOne)

		assert.Error(t, err)
		assert.Equal(t, false, userError)
		assert.Equal(t, true, lowFunds)
		assert.Equal(t, "you don't have enough funds", err.Error())
	})
}

func TestGetFunds(t *testing.T) {
	t.Run("FundsGetOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var testBalanceOneGetLocal = models.Balance{
			UserId: 1,
		}

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGetLocal).DoAndReturn(func(user *models.Balance) (int, error) {
			user.Balance = 1000
			return utils.NO_ERROR, nil
		})

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, err := fundsUseCase.Get(&testBalanceOneGetLocal)

		assert.NoError(t, err)
		assert.Equal(t, false, userError)
	})

	t.Run("FundsAddInvalidUserId", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, err := fundsUseCase.Get(&testBalanceWrongGet)

		assert.Error(t, err)
		assert.Equal(t, true, userError)
	})

	t.Run("DBErrorFirst", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.SERVER_ERROR, errors.New("db error"))

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, err := fundsUseCase.Get(&testBalanceOneGet)

		assert.Error(t, err)
		assert.Equal(t, false, userError)
	})

	t.Run("DBErrorSecond", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceTwoGet).Return(utils.USER_ERROR, errors.New("no user"))
		mockRepoBalance.EXPECT().InsertUser(&testBalanceTwoGet).Return(errors.New("db error"))

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, err := fundsUseCase.Get(&testBalanceTwoGet)

		assert.Error(t, err)
		assert.Equal(t, false, userError)
	})

	t.Run("NewUser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceTwoGet).Return(utils.USER_ERROR, errors.New("no user"))
		mockRepoBalance.EXPECT().InsertUser(&testBalanceTwoGet).Return(nil)

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, err := fundsUseCase.Get(&testBalanceTwoGet)

		assert.NoError(t, err)
		assert.Equal(t, false, userError)
	})
}

func TestTransferFunds(t *testing.T) {
	t.Run("FundsTransferOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		var testBalanceOneGetLocal = models.Balance{
			UserId: 1,
		}

		var testBalanceTwoGetLocal = models.Balance{
			UserId: 2,
		}

		firstMock := mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceTwoGetLocal).DoAndReturn(func(user *models.Balance) (int, error) {
			user.Balance = 10
			return utils.NO_ERROR, nil
		})
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGetLocal).DoAndReturn(func(user *models.Balance) (int, error) {
			user.Balance = 100
			return utils.NO_ERROR, nil
		}).After(firstMock)

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)
		mockRepoTxs.EXPECT().Add(&testTxOneTransfer).Return(nil)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Transfer(&testTxOneTransfer)

		assert.NoError(t, err)
		assert.Equal(t, utils.GetOperationType("Transfer"), testTxOneTransfer.OperationType)
		assert.Equal(t, float64(0), testTxOneTransfer.BalanceFrom)
		assert.Equal(t, float64(110), testTxOneTransfer.Balance)
		assert.Equal(t, false, userError)
		assert.Equal(t, false, lowFunds)
	})

	t.Run("DBErrorFirst", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceTwoGet).Return(utils.SERVER_ERROR, errors.New("db error"))

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Transfer(&testTxOneTransfer)

		assert.Error(t, err)
		assert.Equal(t, false, lowFunds)
		assert.Equal(t, false, userError)
	})

	t.Run("DBErrorSecond", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceTwoGet).Return(utils.USER_ERROR, errors.New("no user"))
		mockRepoBalance.EXPECT().InsertUser(&testBalanceTwoGet).Return(errors.New("db error"))

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Transfer(&testTxOneTransfer)

		assert.Error(t, err)
		assert.Equal(t, false, lowFunds)
		assert.Equal(t, false, userError)
	})

	t.Run("DBErrorThird", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceTwoGet).Return(utils.NO_ERROR, nil)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.SERVER_ERROR, errors.New("db error"))

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Transfer(&testTxOneTransfer)

		assert.Error(t, err)
		assert.Equal(t, false, lowFunds)
		assert.Equal(t, false, userError)
	})

	t.Run("DBErrorFourth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceTwoGet).Return(utils.NO_ERROR, nil)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.USER_ERROR, errors.New("no user"))
		mockRepoBalance.EXPECT().InsertUser(&testBalanceOneGet).Return(errors.New("db error"))

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Transfer(&testTxOneTransfer)

		assert.Error(t, err)
		assert.Equal(t, false, lowFunds)
		assert.Equal(t, false, userError)
	})

	t.Run("WrongSum", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceTwoGet).Return(utils.NO_ERROR, nil)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.NO_ERROR, nil)

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Transfer(&testTxWrongTransfer)

		assert.Error(t, err)
		assert.Equal(t, true, userError)
		assert.Equal(t, false, lowFunds)
		assert.Equal(t, "sum must be positive", err.Error())
	})

	t.Run("LowFunds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var testBalanceOneGetLocal = models.Balance{
			UserId: 1,
		}

		var testBalanceTwoGetLocal = models.Balance{
			UserId: 2,
		}

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceTwoGetLocal).DoAndReturn(func(user *models.Balance) (int, error) {
			user.Balance = 50
			return utils.NO_ERROR, nil
		})
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGetLocal).Return(utils.NO_ERROR, nil)

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, lowFunds, err := fundsUseCase.Transfer(&testTxOneTransfer)

		assert.Error(t, err)
		assert.Equal(t, false, userError)
		assert.Equal(t, true, lowFunds)
		assert.Equal(t, "user doesn't have enough funds", err.Error())
	})
}

func TestGetTransactions(t *testing.T) {
	t.Run("TxsGetOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)

		var testBalanceOneGetLocal = models.Balance{
			UserId: 1,
		}

		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGetLocal).DoAndReturn(func(user *models.Balance) (int, error) {
			user.Balance = 1000
			return utils.NO_ERROR, nil
		})

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)
		mockRepoTxs.EXPECT().GetUserTransactions(&testUserOne, limitInt, since, sort, descBool).Return(testTransactions, utils.NO_ERROR, nil)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, txs, err := fundsUseCase.GetTransactions(&testUserOne, limitInt, since, sort, descBool)

		assert.NoError(t, err)
		assert.Equal(t, testTransactions, txs)
		assert.Equal(t, false, userError)
	})

	t.Run("InvalidUserId", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, txs, err := fundsUseCase.GetTransactions(&testUserWrong, limitInt, since, sort, descBool)

		assert.Error(t, err)
		assert.Equal(t, "incorrect user id", err.Error())
		assert.Equal(t, true, userError)
		assert.Equal(t, []models.Transaction{}, txs)
	})

	t.Run("DBErrorFirst", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.SERVER_ERROR, errors.New("db error"))

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, _, err := fundsUseCase.GetTransactions(&testUserOne, limitInt, since, sort, descBool)

		assert.Error(t, err)
		assert.Equal(t, false, userError)
	})

	t.Run("DBErrorSecond", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.NO_ERROR, nil)

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)
		mockRepoTxs.EXPECT().GetUserTransactions(&testUserOne, limitInt, since, sort, descBool).Return([]models.Transaction{}, utils.SERVER_ERROR, errors.New("db error"))

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, _, err := fundsUseCase.GetTransactions(&testUserOne, limitInt, since, sort, descBool)

		assert.Error(t, err)
		assert.Equal(t, false, userError)
	})

	t.Run("NewUser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.USER_ERROR, errors.New("user doens't exist"))

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, txs, err := fundsUseCase.GetTransactions(&testUserOne, limitInt, since, sort, descBool)

		assert.NoError(t, err)
		assert.Equal(t, []models.Transaction{}, txs)
		assert.Equal(t, false, userError)
	})

	t.Run("UserError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBalance := repository.NewMockBalanceRepoI(ctrl)
		mockRepoBalance.EXPECT().GetBalanceByUserId(&testBalanceOneGet).Return(utils.NO_ERROR, nil)

		mockRepoTxs := repository.NewMockTransactionsRepoI(ctrl)
		mockRepoTxs.EXPECT().GetUserTransactions(&testUserOne, limitInt, since, sort, descBool).Return([]models.Transaction{}, utils.USER_ERROR, errors.New("user error"))

		fundsUseCase := FundsUC{
			BalanceRepo:      mockRepoBalance,
			TransactionsRepo: mockRepoTxs,
		}

		userError, txs, err := fundsUseCase.GetTransactions(&testUserOne, limitInt, since, sort, descBool)

		assert.Error(t, err)
		assert.Equal(t, []models.Transaction{}, txs)
		assert.Equal(t, true, userError)
	})
}
