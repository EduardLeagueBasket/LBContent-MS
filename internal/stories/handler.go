package stories

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateStory(c *gin.Context) {
	var dto CreateStoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido: " + err.Error()})
		return
	}

	if err := h.service.CreateStory(c.Request.Context(), dto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la historia: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": dto})
}

func (h *Handler) GetFeedByPublisher(c *gin.Context) {
	publisherId := c.Param("publisher_id")
	if publisherId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Publisher ID is required"})
		return
	}

	feed, err := h.service.GetFeedByPublisher(c.Request.Context(), uuid.MustParse(publisherId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el feed del publisher: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": feed})
}

func (h *Handler) GetFeedByCountry(c *gin.Context) {
	countryCode := c.Param("country_code")
	if countryCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Country code is required"})
		return
	}

	feed, err := h.service.GetFeedByCountry(c.Request.Context(), countryCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el feed del país: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": feed})
}
