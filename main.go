package main

import (
	"github.com/google/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	balance_handlers "github.com/saskamegaprogrammist/userBalanceService/handlers"
	"github.com/saskamegaprogrammist/userBalanceService/repository"
	"github.com/saskamegaprogrammist/userBalanceService/useCases"
	"github.com/saskamegaprogrammist/userBalanceService/utils"
	"net/http"
	"time"
)

func main() {

	// logger initialization
	utils.LoggerSetup()
	defer utils.LoggerClose()

	// database initialization
	err := repository.Init(pgx.ConnConfig{
		Database: utils.DBName,
		Host:     "localhost",
		User:     "alexis",
		Password: "sinope27",
	})
	if err != nil {
		logger.Fatalf("Couldn't initialize database: %v", err)
	}

	err = useCases.Init(repository.GetBalanceRepo(), repository.GetTransactionsRepo())
	if err != nil {
		logger.Fatalf("Couldn't initialize useCases: %v", err)
	}

	err = balance_handlers.Init(useCases.GetFundsUC())
	if err != nil {
		logger.Fatalf("Couldn't initialize handlers: %v", err)
	}

	// router initialization

	r := mux.NewRouter()
	r.HandleFunc(utils.GetAPIAddress("addFunds"),  balance_handlers.GetUFundsH().Add).Methods("POST")
	r.HandleFunc(utils.GetAPIAddress("withdrawFunds"), balance_handlers.GetUFundsH().Withdraw).Methods("POST")
	r.HandleFunc(utils.GetAPIAddress("getFunds"),  balance_handlers.GetUFundsH().GetBalance).Methods("POST")
	r.HandleFunc(utils.GetAPIAddress("transferFunds"), balance_handlers.GetUFundsH().Transfer).Methods("POST")
	r.HandleFunc(utils.GetAPIAddress("getTransactions"),  balance_handlers.GetUFundsH().GetTransactions).Methods("POST")

	cors := handlers.CORS(handlers.AllowCredentials(), handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"}))

	// server initialization

	server := &http.Server{
		Addr: utils.PortNum,
		Handler : cors(r),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = server.ListenAndServe()

	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}