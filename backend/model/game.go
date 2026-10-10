package model

import (
	"time"
)

type Game struct {
	ID          int        `json:"id"`
	ReleaseDate time.Time  `json:"releaseDate"`
	Categories  []Category `json:"categories"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	PosterUrl   string     `json:"posterUrl"`
}

type GameDetail struct {
	ID                  int        `json:"id"`
	ReleaseDate         time.Time  `json:"releaseDate"`
	Categories          []Category `json:"categories"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Price               float64    `json:"price"`
	PosterUrl           string     `json:"posterUrl"`
	Account             Account    `json:"account"`
	HardwareRequirement []Hardware `json:"hardwareRequirement"`
	Media               []Media    `json:"media"`
}
