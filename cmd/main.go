package main

import (
	"urlShortner/internals/db"
	"urlShortner/internals/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	dbClient := db.NewRedisClient()
	defer dbClient.Close()

	router.SetupRoutes(r, dbClient)

	r.Run(":8000")
}
