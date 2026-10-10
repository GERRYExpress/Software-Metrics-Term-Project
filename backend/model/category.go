package model

type Category struct {
	ID       ID     `json:"id"`
	Name     Name   `json:"name"`
	IsHidden Hidden `json:"isHidden"`
}
