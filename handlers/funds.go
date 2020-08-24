package handlers

import (
	"fmt"
	"github.com/google/logger"
	easy_json "github.com/mailru/easyjson"
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
	err := easy_json.UnmarshalFromReader(req.Body, &newTransaction)
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
	err := easy_json.UnmarshalFromReader(req.Body, &newTransaction)
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
	query := req.URL.Query()
	currency := query.Get("currency")
	var newBalance models.Balance
	var newUserId models.UserId
	err := easy_json.UnmarshalFromReader(req.Body, &newUserId)
	if err != nil {
		jsonError := fmt.Sprintf("Error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(jsonError))
		return
	}
	newBalance.UserId = newUserId.UserId
	newBalance.Currency = utils.CURRENCY
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
	if currency != "" {
		httpClient := &http.Client{}
		address := fmt.Sprintf("%s%s%s", utils.CURRENCY_API, utils.CURRENCY_API_BASE, currency)
		request, err := http.NewRequest("GET", address, nil)
		if err != nil {
			logger.Error(err)
			utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(err.Error()))
			return
		}
		response, err := httpClient.Do(request)
		if err != nil {
			logger.Error(err)
			utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(err.Error()))
			return
		}

		var newCurrency models.CurrencyAll
		err = easy_json.UnmarshalFromReader(response.Body, &newCurrency)
		if err != nil {
			jsonError := fmt.Sprintf("Error unmarshaling json: %v", err.Error())
			logger.Errorf(jsonError)
			utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(jsonError))
			return
		}
		unmarshalledValue, err := newCurrency.GetRatesFieldValueByName(currency)
		if err != nil {
			logger.Errorf(err.Error())
			utils.CreateErrorAnswerJson(writer, utils.StatusCode("Bad Request"), models.CreateMessage(err.Error()))
			return
		}
		newBalance.Balance *= unmarshalledValue
		newBalance.Currency = currency
	}
	utils.CreateAnswerBalanceJson(writer, utils.StatusCode("OK"), newBalance)
}

func (fh *FundsHandlers) Transfer(writer http.ResponseWriter, req *http.Request) {
	var newTransaction models.Transaction
	err := easy_json.UnmarshalFromReader(req.Body, &newTransaction)
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
	err = easy_json.UnmarshalFromReader(req.Body, &newUserId)
	if err != nil {
		jsonError := fmt.Sprintf("Error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(jsonError))
		return
	}
	badRequest, txs, err := fh.FundsUC.GetTransactions(&newUserId, limitInt, since, sort, descBool)
	if badRequest {
		logger.Error(err)
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
