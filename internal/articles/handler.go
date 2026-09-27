package articles

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateArticleHandler(c *gin.Context) {
	var dto CreateArticleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido: " + err.Error()})
		return
	}

	if err := h.service.CreateArticle(c.Request.Context(), dto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el artículo: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artículo creado exitosamente", "data": dto})
}

func (h *Handler) ListArticlesHandler(c *gin.Context) {
	countryCode := c.Query("country_code")
	entityType := c.Query("entity_type")
	entityID := c.Query("entity_id")

	var (
		items []ArticleDTO
		err   error
	)
	switch {
	case countryCode != "":
		items, err = h.service.ListByCountry(c.Request.Context(), countryCode)
	case entityType != "" && entityID != "":
		items, err = h.service.ListByEntity(c.Request.Context(), ListByEntityDTO{
			EntityType: entityType,
			EntityID:   entityID,
		})
	default:
		items, err = h.service.ListArticles(c.Request.Context())
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar los artículos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) GetArticleHandler(c *gin.Context) {
	article, err := h.service.GetArticle(c.Request.Context(), c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artículo no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el artículo: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *Handler) UpdateArticleHandler(c *gin.Context) {
	var dto ArticleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido: " + err.Error()})
		return
	}

	if err := h.service.UpdateArticle(c.Request.Context(), c.Param("id"), dto); err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artículo no encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el artículo: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artículo actualizado exitosamente", "data": dto})
}

func (h *Handler) DeleteArticleHandler(c *gin.Context) {
	if err := h.service.DeleteArticleByID(c.Request.Context(), c.Param("id")); err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artículo no encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el artículo: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artículo eliminado exitosamente"})
}
