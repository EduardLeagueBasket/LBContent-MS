package article_categories

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	articleCategoryGroup := r.Group("/article-categories")
	{
		articleCategoryGroup.POST("", h.CreateArticleCategoryHandler)
		articleCategoryGroup.GET("", h.ListArticleCategoriesHandler)
	}
}
