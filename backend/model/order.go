package model

type Order struct {
	ID            ID            `json:"id"`
	UserID        ID            `json:"userId"`
	OrderDate     OrderDate     `json:"orderDate"`
	PaymentMethod PaymentMethod `json:"paymentMethod"`
	FinalPrice    FinalPrice    `json:"finalPrice"`
	Games         []Game        `json:"games"`
}
