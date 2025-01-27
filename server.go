package server

import (
	"myfavouritemovies/routes"

	"github.com/gin-gonic/gin"
)

func CreateServer() *gin.Engine {
	router := gin.Default()
	routes.SetUpRoutes(router)
	return router
}
