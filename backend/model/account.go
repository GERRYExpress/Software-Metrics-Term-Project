package model

type Account struct {
	ID          ID       `json:"id"`
	Username    Username `json:"username"`
	Password    Password `json:"-"`
	CreatedDate Date     `json:"created_date"`
	IsBuyer     Hidden   `json:"is_buyer"`
	IsSeller    Hidden   `json:"is_seller"`
	TotalGame   Count    `json:"total_game"`
	TotalSpend  Amount   `json:"total_spend"`
	TotalSale   Amount   `json:"total_sale"`
}
