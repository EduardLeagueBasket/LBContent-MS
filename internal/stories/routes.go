package stories

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	storiesGroup := r.Group("/stories")
	{
		storiesGroup.POST("", h.CreateStory)
		storiesGroup.GET("/publisher/:publisher_id", h.GetFeedByPublisher)
		storiesGroup.GET("/country/:country_code", h.GetFeedByCountry)
	}
}
