package model

import "time"

type Cart struct {
	UserID    int       `json:"userId"`
	ProductID int       `json:"productId"`
	CartDate  time.Time `json:"cartDate"`
	Game      Game      `json:"game"`
}
