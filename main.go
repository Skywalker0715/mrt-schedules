package main

import (
	"github.com/Skywalker0715/mrt-schedules/modules/station"
	"github.com/gin-gonic/gin"
)

func main() {
	initializeRoutes()
}

func initializeRoutes() {
	router := gin.Default()
	api := router.Group("/v1/api")

	// Initialize station routes
	station.InitializeRouter(api)

	router.Run(":8080")
}
