package publisher

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

func (h *Handler) UpsertHandler(c *gin.Context) {
	var dto UpdatePublisherDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido: " + err.Error()})
		return
	}

	if err := h.service.UpsertPublisher(c.Request.Context(), dto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el publisher: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Publisher procesado exitosamente",
		"data":    dto,
	})
}

func (h *Handler) ListHandler(c *gin.Context) {
	dto := ListPublishersDTO{EntityType: c.Query("entity_type")}
	items, err := h.service.ListPublishers(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) CreatePlatformHandler(c *gin.Context) {
	var dto CreatePlatformPublisherDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido: " + err.Error()})
		return
	}
	created, err := h.service.CreatePlatformPublisher(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}
