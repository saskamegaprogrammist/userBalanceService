package utils

var statusCodes = map[string]int {
	"OK" :          200,
	"Created" :     201,
	"Bad Request" : 400,
	"Unauthorized" : 401,
	"Payment Required" :402,
	"Not Found" : 404,
	"Method Not Allowed" : 405,
	"Conflict" : 409,
	"Internal Server Error": 500,
	"Not Implemented": 501,
	"Bad Gateway" : 502,
	"Service Unavailable" : 503,
}

var API = map[string]string {
	"addFunds" : "/funds/add",
	"withdrawFunds" : "/funds/withdraw",
	"getFunds" : "/funds/get",
	"transferFunds" : "/funds/transfer",
	"getTransactions": "/funds/details",
}

func StatusCode(mess string) int {
	return statusCodes[mess]
}

func GetAPIAddress(address string) string {
	return API[address]
}

const ERROR_ID  = 0
const LIMIT_DEFAULT = -1
const LogFile  = "log.log"
const DBName = "user_balance_service"
const PortNum = ":5000"

const (
	NO_ERROR = iota
	USER_ERROR
	SERVER_ERROR
)

var operationTypes = map[string]int {
	"Add" : 1,
	"Withdraw" : 2,
	"Transfer" : 3,
}

func GetOperationType(operation string) int {
	return operationTypes[operation]
}

