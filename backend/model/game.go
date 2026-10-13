package model

type Game struct {
	ID          ID          `json:"id"`
	UserID      ID          `json:"userId"`
	ReleaseDate ReleaseDate `json:"releaseDate"`
	Categories  Categories  `json:"categories"`
	Name        Name        `json:"name"`
	Description Description `json:"description"`
	Price       Price       `json:"price"`
	PosterUrl   PosterURL   `json:"posterUrl"`
}

type GameDetail struct {
	ID                  ID          `json:"id"`
	ReleaseDate         ReleaseDate `json:"releaseDate"`
	Categories          Categories  `json:"categories"`
	Name                Name        `json:"name"`
	Description         Description `json:"description"`
	Price               Price       `json:"price"`
	PosterUrl           PosterURL   `json:"posterUrl"`
	Account             Account     `json:"account"`
	HardwareRequirement []Hardware  `json:"hardwareRequirement"`
	Media               []Media     `json:"media"`
}

type Categories []Category

func (g *Game) AddCategory(category Category) {
	g.Categories = append(g.Categories, category)
}
