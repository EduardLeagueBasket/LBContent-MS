package publisher

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	publishersGroup := r.Group("/publishers")
	{
		publishersGroup.GET("", h.ListHandler)
		publishersGroup.POST("", h.CreatePlatformHandler)
		publishersGroup.POST("/sync", h.UpsertHandler)
	}
}
