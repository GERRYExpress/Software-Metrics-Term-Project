package handler

import (
	"errors"
	"fmt"
	"net/http"
	"software-metrics-project/model"
	"software-metrics-project/repository"
	"software-metrics-project/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	Service *service.GameService
}

type GameCreateRequest struct {
	UserID               int                     `json:"userId"`
	ReleaseDate          model.ReleaseDate       `json:"releaseDate"`
	CategoryIDs          []int                   `json:"categoryIds"`
	Name                 string                  `json:"name"`
	Description          string                  `json:"description"`
	Price                float64                 `json:"price"`
	PosterURL            string                  `json:"posterUrl"`
	HardwareRequirements []HardwareCreateRequest `json:"hardwareRequirements"`
	MediaURLs            []string                `json:"mediaUrls"`
}

type HardwareCreateRequest struct {
	OS             string `json:"operatingSystem"`
	Processor      string `json:"processor"`
	Memory         string `json:"memory"`
	GraphicsCard   string `json:"graphicsCard"`
	DirectXVersion string `json:"directXVersion"`
	Network        string `json:"network"`
	Storage        string `json:"storage"`
	Soundcard      string `json:"soundcard"`
	AdditionalNote string `json:"additionalNote"`
	Type           string `json:"type"`
}

type GameResponse struct {
	ID                  model.ID          `json:"id"`
	ReleaseDate         model.ReleaseDate `json:"releaseDate"`
	Categories          []model.Category  `json:"categories"`
	Name                model.Name        `json:"name"`
	Description         model.Description `json:"description"`
	Price               model.Price       `json:"price"`
	PosterURL           model.PosterURL   `json:"posterUrl"`
	Account             model.Account     `json:"account"`
	HardwareRequirement []model.Hardware  `json:"hardwareRequirement"`
	Media               []model.Media     `json:"media"`
}

func (h *GameHandler) Create(c *gin.Context) {
	var request GameCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.UserID < 1 || strings.TrimSpace(request.Name) == "" ||
		strings.TrimSpace(request.Description) == "" ||
		request.Price < 0 || strings.TrimSpace(request.PosterURL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game details"})
		return
	}
	created, err := h.Service.Create(request.toModel())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newGameResponse(*created))
}

func (h *GameHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}
	var request GameCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.UserID < 1 || strings.TrimSpace(request.Name) == "" ||
		strings.TrimSpace(request.Description) == "" ||
		request.Price < 0 || strings.TrimSpace(request.PosterURL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game details"})
		return
	}
	updated, err := h.Service.Update(id, request.toModel())
	if err != nil {
		if errors.Is(err, repository.ErrGameNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newGameResponse(*updated))
}

func (h *GameHandler) Delete(c *gin.Context) {
	id, err := positiveID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, err := positiveID(c.Query("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId must be a positive integer"})
		return
	}
	if err := h.Service.Delete(id, userID); err != nil {
		switch {
		case errors.Is(err, repository.ErrGameNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		case errors.Is(err, repository.ErrGameNotOwned):
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this game"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (r GameCreateRequest) toModel() model.GameDetail {
	categories := make([]model.Category, 0, len(r.CategoryIDs))
	for _, id := range r.CategoryIDs {
		categories = append(categories, model.NewCategory(id, "", false))
	}
	requirements := make([]model.Hardware, 0, len(r.HardwareRequirements))
	for _, requirement := range r.HardwareRequirements {
		requirements = append(requirements, model.NewHardware(model.HardwareData{
			OS: requirement.OS, Processor: requirement.Processor, Memory: requirement.Memory,
			GraphicsCard: requirement.GraphicsCard, DirectXVersion: requirement.DirectXVersion,
			Network: requirement.Network, Storage: requirement.Storage, Soundcard: requirement.Soundcard,
			AdditionalNote: requirement.AdditionalNote, Type: requirement.Type,
		}))
	}
	media := make([]model.Media, 0, len(r.MediaURLs))
	for _, url := range r.MediaURLs {
		media = append(media, model.NewMedia(0, url))
	}
	return model.GameDetail{
		ReleaseDate: r.ReleaseDate,
		Categories:  categories,
		Name:        model.NewName(r.Name), Description: model.NewDescription(r.Description),
		Price: model.NewPrice(r.Price), PosterUrl: model.NewPosterURL(r.PosterURL),
		Account:             model.NewAccount(r.UserID, ""),
		HardwareRequirement: requirements, Media: media,
	}
}

func newGameResponse(game model.GameDetail) GameResponse {
	return GameResponse{
		ID: game.ID, ReleaseDate: game.ReleaseDate, Categories: game.Categories,
		Name: game.Name, Description: game.Description, Price: game.Price,
		PosterURL: game.PosterUrl, Account: game.Account,
		HardwareRequirement: game.HardwareRequirement, Media: game.Media,
	}
}

func (h *GameHandler) GetGames(c *gin.Context) {
	var (
		games []model.Game
		err   error
	)
	limit, page, err := pagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	keyword := c.Query("keyword")
	category := c.Query("category")
	var sellerID *int
	if rawSellerID := c.Query("sellerId"); rawSellerID != "" {
		parsedSellerID, parseErr := strconv.Atoi(rawSellerID)
		if parseErr != nil || parsedSellerID < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "sellerId must be a positive integer"})
			return
		}
		sellerID = &parsedSellerID
	}
	if keyword != "" {
		games, err = h.Service.GetGamesByName(keyword, category, sellerID, limit, page)
	} else {
		games, err = h.Service.GetGames(category, sellerID, limit, page)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	totalCount, err := h.Service.GetGamesCount(keyword, category, sellerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("X-Total-Count", strconv.Itoa(totalCount))
	c.JSON(http.StatusOK, games)
}

func pagination(c *gin.Context) (int, int, error) {
	limit, err := strconv.Atoi(c.DefaultQuery("_limit", "10"))
	if err != nil || limit < 1 {
		return 0, 0, fmt.Errorf("_limit must be a positive integer")
	}

	page, err := strconv.Atoi(c.DefaultQuery("_page", "1"))
	if err != nil || page < 1 {
		return 0, 0, fmt.Errorf("_page must be a positive integer")
	}

	return limit, page, nil
}

func (h *GameHandler) GetGame(c *gin.Context) {
	gameID := c.Param("id")

	id, err := strconv.Atoi(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := h.Service.GetGame(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, game)
}
