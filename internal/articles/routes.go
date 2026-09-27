package articles

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	articleGroup := r.Group("/articles")
	{
		articleGroup.POST("", h.CreateArticleHandler)
		articleGroup.GET("", h.ListArticlesHandler)
		articleGroup.GET("/:id", h.GetArticleHandler)
		articleGroup.PUT("/:id", h.UpdateArticleHandler)
		articleGroup.DELETE("/:id", h.DeleteArticleHandler)
	}
}
