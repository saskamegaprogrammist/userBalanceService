package handlers

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/saskamegaprogrammist/userBalanceService/models"
	"github.com/saskamegaprogrammist/userBalanceService/useCases"
	"github.com/saskamegaprogrammist/userBalanceService/utils"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
	"time"
)

var fh FundsHandlers

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

var testUserOne = models.UserId{
	UserId: 1,
}

var testUserWrong = models.UserId{
	UserId: 0,
}

var testBalanceOneGet = models.Balance{
	UserId:   1,
	Currency: utils.CURRENCY,
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

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Add(&testTxOne).Return(false, nil)

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_id": %v, "sum": %v}`, testTxOne.UserId, testTxOne.Sum)

		apitest.New("FundsAddOK").
			Handler(http.HandlerFunc(fh.Add)).
			Method("Post").
			URL(utils.GetAPIAddress("addFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("UserIdWrong", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Add(&testTxWrong).Return(true, errors.New("invalid user id"))

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_id": %v, "sum": %v}`, testTxWrong.UserId, testTxWrong.Sum)

		apitest.New("UserIdWrong").
			Handler(http.HandlerFunc(fh.Add)).
			Method("Post").
			URL(utils.GetAPIAddress("addFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("InternalError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Add(&testTxOne).Return(false, errors.New("db error"))

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_id": %v, "sum": %v}`, testTxOne.UserId, testTxOne.Sum)

		apitest.New("InternalError").
			Handler(http.HandlerFunc(fh.Add)).
			Method("Post").
			URL(utils.GetAPIAddress("addFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()

	})

	t.Run("MalformedJSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_i: %v, "sum": %v}`, testTxOne.UserId, testTxOne.Sum)

		apitest.New("MalformedJSON").
			Handler(http.HandlerFunc(fh.Add)).
			Method("Post").
			URL(utils.GetAPIAddress("addFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusInternalServerError).
			Assert(jsonpath.Contains(`$.message`, "Error unmarshaling json")).
			End()
	})
}

func TestWithdrawFunds(t *testing.T) {
	t.Run("FundsWithdrawOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Withdraw(&testTxTwo).Return(false, false, nil)

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_id": %v, "sum": %v}`, testTxTwo.UserId, testTxTwo.Sum)

		apitest.New("FundsAddOK").
			Handler(http.HandlerFunc(fh.Withdraw)).
			Method("Post").
			URL(utils.GetAPIAddress("withdrawFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("UserIdWrong", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Withdraw(&testTxWrong).Return(true, false, errors.New("invalid user id"))

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_id": %v, "sum": %v}`, testTxWrong.UserId, testTxWrong.Sum)

		apitest.New("UserIdWrong").
			Handler(http.HandlerFunc(fh.Withdraw)).
			Method("Post").
			URL(utils.GetAPIAddress("withdrawFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("UserLowFunds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Withdraw(&testTxTwo).Return(false, true, errors.New("low funds error"))

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_id": %v, "sum": %v}`, testTxTwo.UserId, testTxTwo.Sum)

		apitest.New("UserLowFunds").
			Handler(http.HandlerFunc(fh.Withdraw)).
			Method("Post").
			URL(utils.GetAPIAddress("withdrawFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusPaymentRequired).
			End()
	})

	t.Run("InternalError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Withdraw(&testTxTwo).Return(false, false, errors.New("db error"))

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_id": %v, "sum": %v}`, testTxTwo.UserId, testTxTwo.Sum)

		apitest.New("InternalError").
			Handler(http.HandlerFunc(fh.Withdraw)).
			Method("Post").
			URL(utils.GetAPIAddress("withdrawFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()

	})

	t.Run("MalformedJSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"use_i: %v, "sum": %v}`, testTxOne.UserId, testTxOne.Sum)

		apitest.New("MalformedJSON").
			Handler(http.HandlerFunc(fh.Withdraw)).
			Method("Post").
			URL(utils.GetAPIAddress("withdrawFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusInternalServerError).
			Assert(jsonpath.Contains(`$.message`, "Error unmarshaling json")).
			End()
	})
}

func TestGetFunds(t *testing.T) {
	t.Run("FundsGetSimpleOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Get(&testBalanceOneGet).Return(false, nil)

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user": %v}`, testUserOne.UserId)

		apitest.New("FundsGetSimpleOK").
			Handler(http.HandlerFunc(fh.GetBalance)).
			Method("Post").
			URL(utils.GetAPIAddress("getFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Matches("$.balance", `([0-9]*[.])?[0-9]+`)).
			End()
	})

	t.Run("UserIdError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Get(&testBalanceWrongGet).Return(true, errors.New("user id error"))

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user": %v}`, testUserWrong.UserId)

		apitest.New("UserIdError").
			Handler(http.HandlerFunc(fh.GetBalance)).
			Method("Post").
			URL(utils.GetAPIAddress("getFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("InternalError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Get(&testBalanceOneGet).Return(false, errors.New("db error"))

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user": %v}`, testUserOne.UserId)

		apitest.New("UserIdError").
			Handler(http.HandlerFunc(fh.GetBalance)).
			Method("Post").
			URL(utils.GetAPIAddress("getFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_id %v}`, testUserOne.UserId)

		apitest.New("MalformedJSON").
			Handler(http.HandlerFunc(fh.GetBalance)).
			Method("Post").
			URL(utils.GetAPIAddress("getFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusInternalServerError).
			Assert(jsonpath.Contains(`$.message`, "Error unmarshaling json")).
			End()
	})

	t.Run("FundsGetCurrencyOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Get(&testBalanceOneGet).Return(false, nil)

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user": %v}`, testUserOne.UserId)

		apitest.New("FundsGetSimpleOK").
			Handler(http.HandlerFunc(fh.GetBalance)).
			Method("Post").
			URL(utils.GetAPIAddress("getFunds")).
			Query("currency", "EUR").
			Body(jsonBody).
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Matches("$.balance", `([0-9]*[.])?[0-9]+`)).
			Assert(jsonpath.Matches("$.currency", `EUR`)).
			End()
	})

	t.Run("FundsGetCurrencyWrong", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Get(&testBalanceOneGet).Return(false, nil)

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user": %v}`, testUserOne.UserId)

		apitest.New("FundsGetSimpleOK").
			Handler(http.HandlerFunc(fh.GetBalance)).
			Method("Post").
			URL(utils.GetAPIAddress("getFunds")).
			Query("currency", "dsgsdg").
			Body(jsonBody).
			Expect(t).
			Status(http.StatusBadRequest).
			Assert(jsonpath.Contains("$.message", "invalid currency")).
			End()
	})
}

func TestTransferFunds(t *testing.T) {
	t.Run("FundsTransferOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Transfer(&testTxOneTransfer).Return(false, false, nil)

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_id": %v, "user_from_id": %v, "sum": %v}`, testTxOneTransfer.UserId, testTxOneTransfer.UserFromId, testTxOneTransfer.Sum)

		apitest.New("FundsGetSimpleOK").
			Handler(http.HandlerFunc(fh.Transfer)).
			Method("Post").
			URL(utils.GetAPIAddress("transferFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("FundsTransferWrong", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Transfer(&testTxWrongTransfer).Return(true, false, errors.New("wrong sum"))

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_id": %v, "user_from_id": %v, "sum": %v}`, testTxWrongTransfer.UserId, testTxWrongTransfer.UserFromId, testTxWrongTransfer.Sum)

		apitest.New("FundsTransferWrong").
			Handler(http.HandlerFunc(fh.Transfer)).
			Method("Post").
			URL(utils.GetAPIAddress("transferFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("FundsTransferLowFunds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Transfer(&testTxOneTransfer).Return(false, true, errors.New("low funds"))

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_id": %v, "user_from_id": %v, "sum": %v}`, testTxOneTransfer.UserId, testTxOneTransfer.UserFromId, testTxOneTransfer.Sum)

		apitest.New("FundsTransferLowFunds").
			Handler(http.HandlerFunc(fh.Transfer)).
			Method("Post").
			URL(utils.GetAPIAddress("transferFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusPaymentRequired).
			End()
	})

	t.Run("InternalError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().Transfer(&testTxOneTransfer).Return(false, false, errors.New("db error"))

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_id": %v, "user_from_id": %v, "sum": %v}`, testTxOneTransfer.UserId, testTxOneTransfer.UserFromId, testTxOneTransfer.Sum)

		apitest.New("FundsTransferLowFunds").
			Handler(http.HandlerFunc(fh.Transfer)).
			Method("Post").
			URL(utils.GetAPIAddress("transferFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user_ %v, "user_from_id": %v, "sum": %v}`, testTxOneTransfer.UserId, testTxOneTransfer.UserFromId, testTxOneTransfer.Sum)

		apitest.New("MalformedJSON").
			Handler(http.HandlerFunc(fh.Transfer)).
			Method("Post").
			URL(utils.GetAPIAddress("transferFunds")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusInternalServerError).
			Assert(jsonpath.Contains(`$.message`, "Error unmarshaling json")).
			End()
	})
}

func TestGetTransactions(t *testing.T) {
	t.Run("TxsGetSimpleOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().GetTransactions(&testUserOne, limitInt, since, sort, descBool).Return(false, testTransactions, nil)

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user": %v}`, testUserOne.UserId)

		apitest.New("TxsGetSimpleOK").
			Handler(http.HandlerFunc(fh.GetTransactions)).
			Method("Post").
			URL(utils.GetAPIAddress("getTransactions")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("TxsGetSimpleWrong", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().GetTransactions(&testUserWrong, limitInt, since, sort, descBool).Return(true, []models.Transaction{}, errors.New("user error"))

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user": %v}`, testUserWrong.UserId)

		apitest.New("TxsGetSimpleWrong").
			Handler(http.HandlerFunc(fh.GetTransactions)).
			Method("Post").
			URL(utils.GetAPIAddress("getTransactions")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("InternalError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().GetTransactions(&testUserOne, limitInt, since, sort, descBool).Return(false, []models.Transaction{}, errors.New("db error"))

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user": %v}`, testUserOne.UserId)

		apitest.New("InternalError").
			Handler(http.HandlerFunc(fh.GetTransactions)).
			Method("Post").
			URL(utils.GetAPIAddress("getTransactions")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"us %v}`, testUserOne.UserId)

		apitest.New("MalformedJSON").
			Handler(http.HandlerFunc(fh.GetTransactions)).
			Method("Post").
			URL(utils.GetAPIAddress("getTransactions")).
			Body(jsonBody).
			Expect(t).
			Status(http.StatusInternalServerError).
			Assert(jsonpath.Contains(`$.message`, "Error unmarshaling json")).
			End()
	})

	t.Run("TxsGetSinceOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().GetTransactions(&testUserOne, limitInt, "2020-08-22T15:04:05.999999-07:00", sort, descBool).Return(false, testTransactions, nil)

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user": %v}`, testUserOne.UserId)

		apitest.New("TxsGetSimpleOK").
			Handler(http.HandlerFunc(fh.GetTransactions)).
			Method("Post").
			URL(utils.GetAPIAddress("getTransactions")).
			Query("since", "2020-08-22T15:04:05.999999-07:00").
			Body(jsonBody).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("TxsGetSinceLimitOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().GetTransactions(&testUserOne, 2, "2020-08-22T15:04:05.999999-07:00", sort, descBool).Return(false, testTransactions, nil)

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user": %v}`, testUserOne.UserId)

		apitest.New("TxsGetSimpleOK").
			Handler(http.HandlerFunc(fh.GetTransactions)).
			Method("Post").
			URL(utils.GetAPIAddress("getTransactions")).
			Query("since", "2020-08-22T15:04:05.999999-07:00").
			Query("limit", "2").
			Body(jsonBody).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("TxsGetWholeOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)
		mockUseCase.EXPECT().GetTransactions(&testUserOne, 2, "2020-08-22T15:04:05.999999-07:00", "sum", true).Return(false, testTransactions, nil)

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user": %v}`, testUserOne.UserId)

		apitest.New("TxsGetSimpleOK").
			Handler(http.HandlerFunc(fh.GetTransactions)).
			Method("Post").
			URL(utils.GetAPIAddress("getTransactions")).
			Query("since", "2020-08-22T15:04:05.999999-07:00").
			Query("limit", "2").
			Query("sort", "sum").
			Query("desc", "true").
			Body(jsonBody).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("TxsGetLimitWrong", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := useCases.NewMockFundsUCInterface(ctrl)

		fh.FundsUC = mockUseCase

		jsonBody := fmt.Sprintf(`{"user": %v}`, testUserOne.UserId)

		apitest.New("TxsGetSimpleOK").
			Handler(http.HandlerFunc(fh.GetTransactions)).
			Method("Post").
			URL(utils.GetAPIAddress("getTransactions")).
			Query("limit", "hello").
			Body(jsonBody).
			Expect(t).
			Status(http.StatusBadRequest).
			Assert(jsonpath.Contains(`$.message`, "bad limit query param")).
			End()
	})
}
