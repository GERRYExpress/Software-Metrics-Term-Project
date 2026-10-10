package model

type Order struct {
	ID            ID            `json:"id"`
	OrderDate     OrderDate     `json:"orderDate"`
	PaymentMethod PaymentMethod `json:"paymentMethod"`
	FinalPrice    FinalPrice    `json:"finalPrice"`
}

type OrderItem struct {
	OrderID ID   `json:"orderId"`
	GameID  ID   `json:"gameId"`
	Game    Game `json:"game"`
}
