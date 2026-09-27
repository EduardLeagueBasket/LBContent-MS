package router

import (
	"github.com/eduardeddy/content-ms/internal/article_categories"
	"github.com/eduardeddy/content-ms/internal/articles"
	"github.com/eduardeddy/content-ms/internal/publisher"
	"github.com/eduardeddy/content-ms/internal/stories"
	"github.com/gin-gonic/gin"
)

func Setup(
	publisherHandler *publisher.Handler,
	storiesHandler *stories.Handler,
	articlesHandler *articles.Handler,
	articleCategoriesHandler *article_categories.Handler,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		publisherHandler.RegisterRoutes(api)
		storiesHandler.RegisterRoutes(api)
		articlesHandler.RegisterRoutes(api)
		articleCategoriesHandler.RegisterRoutes(api)
	}

	return r
}
