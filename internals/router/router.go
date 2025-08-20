package router

import (
	"urlShortner/internals/handlers"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(r *gin.Engine, db *redis.Client) {
	api := r.Group("/api/v1")

	api.GET("/:shortURL", handlers.GetURL(db))
	api.POST("/createURL", handlers.CreateURL(db))
}
