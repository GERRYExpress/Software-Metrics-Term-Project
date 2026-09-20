package model

import "time"

type Account struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	CreatedDate time.Time `json:"created_date"`
	IsBuyer     bool      `json:"is_buyer"`
	IsSeller    bool      `json:"is_seller"`
	TotalGame   int       `json:"total_game"`
	TotalSpend  float64   `json:"total_spend"`
	TotalSale   float64   `json:"total_sale"`
}
