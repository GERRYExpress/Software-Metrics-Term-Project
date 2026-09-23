package handler

import (
	"net/http"
	"software-metrics-project/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	Service *service.CategoryService
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.Service.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}
