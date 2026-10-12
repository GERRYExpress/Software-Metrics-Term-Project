package model

type Account struct {
	ID          ID       `json:"id"`
	Username    Username `json:"username"`
	Password    Password `json:"-"`
	CreatedDate Date     `json:"createdDate"`
	IsBuyer     Hidden   `json:"isBuyer"`
	IsSeller    Hidden   `json:"isSeller"`
	TotalGame   Count    `json:"totalGame"`
	TotalSpend  Amount   `json:"totalSpend"`
	TotalSale   Amount   `json:"totalSale"`
}
