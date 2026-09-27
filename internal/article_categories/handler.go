package article_categories

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateArticleCategoryHandler(c *gin.Context) {
	var dto ArticleCategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido: " + err.Error()})
		return
	}

	created, err := h.service.CreateArticleCategory(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la categoría de artículo: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categoría de artículo creada exitosamente", "data": created})
}

func (h *Handler) ListArticleCategoriesHandler(c *gin.Context) {
	items, err := h.service.ListArticleCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar las categorías de artículo: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}
