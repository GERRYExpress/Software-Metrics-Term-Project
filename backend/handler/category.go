package handler

import (
	"net/http"
	"software-metrics-project/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	Service *service.CategoryService
}

type CategoryRequest struct {
	Name string `json:"name"`
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.Service.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var category CategoryRequest
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := h.Service.Create(category.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}
