package routes

import (
	"myfavouritemovies/routes/handlers"
	"myfavouritemovies/utils"

	"github.com/gin-gonic/gin"
)


func SetUpRoutes(router *gin.Engine) {

	router.POST("/users",handlers.AddUser)

	auth := router.Group("/")
	auth.Use(utils.HardcodedUserMiddleware())
	auth.PATCH("/users/", handlers.UpdateUser)
	auth.POST("/favgenres", handlers.AddFavoriteGenre)
	auth.DELETE("/favgenres", handlers.DeleteFavoriteGenre)
	auth.POST("/movies", handlers.AddFavoriteMovie)
	auth.DELETE("/movies", handlers.DeleteFavoriteMovie)
	auth.PATCH("/movies/toggle", handlers.ToggleWatchedStatus)
	auth.POST("/genres", handlers.AddGenres)
}
