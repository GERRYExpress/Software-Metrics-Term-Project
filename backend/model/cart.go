package model

type Cart struct {
	UserID    ID   `json:"userId"`
	ProductID ID   `json:"productId"`
	CartDate  Date `json:"cartDate"`
	Game      Game `json:"game"`
}
