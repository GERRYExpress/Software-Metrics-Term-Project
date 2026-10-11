package model

type Media struct {
	ID       ID       `json:"id"`
	MediaUrl MediaURL `json:"mediaUrl"`
	Type     string   `json:"type"`
}
