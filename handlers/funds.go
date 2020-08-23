package handlers

import (
	"fmt"
	"github.com/google/logger"
	json "github.com/mailru/easyjson"
	"github.com/saskamegaprogrammist/userBalanceService/models"
	"github.com/saskamegaprogrammist/userBalanceService/useCases"
	"github.com/saskamegaprogrammist/userBalanceService/utils"
	"net/http"
	"strconv"
)

type FundsHandlers struct {
	FundsUC useCases.FundsUCInterface
}

func (fh *FundsHandlers) Add(writer http.ResponseWriter, req *http.Request) {
	var newTransaction models.Transaction
	err := json.UnmarshalFromReader(req.Body, &newTransaction)
	if err != nil {
		jsonError := fmt.Sprintf("Error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(jsonError))
		return
	}
	badRequest, err := fh.FundsUC.Add(&newTransaction)
	if badRequest {
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Bad Request"), models.CreateMessage(err.Error()))
		return
	}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(err.Error()))
		return
	}

	utils.CreateEmptyBodyAnswerJson(writer, utils.StatusCode("OK"))
}

func (fh *FundsHandlers) Withdraw(writer http.ResponseWriter, req *http.Request) {
	var newTransaction models.Transaction
	err := json.UnmarshalFromReader(req.Body, &newTransaction)
	if err != nil {
		jsonError := fmt.Sprintf("Error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(jsonError))
		return
	}
	badRequest, lowFunds, err := fh.FundsUC.Withdraw(&newTransaction)
	if badRequest {
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Bad Request"), models.CreateMessage(err.Error()))
		return
	}
	if lowFunds {
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Payment Required"), models.CreateMessage(err.Error()))
		return
	}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(err.Error()))
		return
	}
	utils.CreateEmptyBodyAnswerJson(writer, utils.StatusCode("OK"))
}

func (fh *FundsHandlers) GetBalance(writer http.ResponseWriter, req *http.Request) {
	var newBalance models.Balance
	var newUserId models.UserId
	err := json.UnmarshalFromReader(req.Body, &newUserId)
	if err != nil {
		jsonError := fmt.Sprintf("Error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(jsonError))
		return
	}
	newBalance.UserId = newUserId.UserId
	badRequest, err := fh.FundsUC.Get(&newBalance)
	if badRequest {
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Bad Request"), models.CreateMessage(err.Error()))
		return
	}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(err.Error()))
		return
	}
	utils.CreateAnswerBalanceJson(writer, utils.StatusCode("OK"), newBalance)
}

func (fh *FundsHandlers) Transfer(writer http.ResponseWriter, req *http.Request) {
	var newTransaction models.Transaction
	err := json.UnmarshalFromReader(req.Body, &newTransaction)
	if err != nil {
		jsonError := fmt.Sprintf("Error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(jsonError))
		return
	}
	badRequest, lowFunds, err := fh.FundsUC.Transfer(&newTransaction)
	if badRequest {
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Bad Request"), models.CreateMessage(err.Error()))
		return
	}
	if lowFunds {
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Payment Required"), models.CreateMessage(err.Error()))
		return
	}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(err.Error()))
		return
	}
	utils.CreateEmptyBodyAnswerJson(writer, utils.StatusCode("OK"))
}

func (fh *FundsHandlers) GetTransactions(writer http.ResponseWriter, req *http.Request) {
	var err error
	query := req.URL.Query()
	limit := query.Get("limit")
	since := query.Get("since")
	desc := query.Get("desc")
	sort := query.Get("sort")
	limitInt := utils.LIMIT_DEFAULT
	descBool := false
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			utils.CreateErrorAnswerJson(writer, utils.StatusCode("Bad Request"), models.CreateMessage("bad limit query param"))
			return
		}
	}
	if desc == "true" {
		descBool = true
	}
	var newUserId models.UserId
	err = json.UnmarshalFromReader(req.Body, &newUserId)
	if err != nil {
		jsonError := fmt.Sprintf("Error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(jsonError))
		return
	}
	badRequest, txs, err := fh.FundsUC.GetTransactions(&newUserId, limitInt, since, sort, descBool)
	if badRequest {
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Bad Request"), models.CreateMessage(err.Error()))
		return
	}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(err.Error()))
		return
	}
	utils.CreateAnswerTransactionsJson(writer, utils.StatusCode("OK"), txs)
}


