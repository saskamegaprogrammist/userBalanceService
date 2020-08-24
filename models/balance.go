package models

//easyjson:json
type Balance struct {
	Id       int     `json:"-"`
	UserId   int     `json:"user_id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}
