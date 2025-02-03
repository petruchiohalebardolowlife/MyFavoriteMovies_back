package routes

import (
	"myfavouritemovies/routes/handlers"
	"myfavouritemovies/utils"

	"github.com/gin-gonic/gin"
)


func SetUpRoutes(router *gin.Engine) {

	router.POST("/users",handlers.AddUserHandler)
  
	auth := router.Group("/")
	auth.Use(utils.HardcodedUserMiddleware())
	auth.PATCH("/users/", handlers.UpdateUserHandler)
	auth.POST("/favgenres", handlers.AddFavoriteGenreHandler)
	auth.DELETE("/favgenres", handlers.DeleteFavoriteGenreHandler)
	auth.POST("/movies", handlers.AddFavoriteMovieHandler)
	auth.DELETE("/movies", handlers.DeleteFavoriteMovieHandler)
	auth.PATCH("/movies/toggle", handlers.ToggleWatchedStatusHandler)
	auth.POST("/genres", handlers.AddGenresHandler)
}
