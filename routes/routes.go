package routes

import (
	"myfavouritemovies/routes/apihandlers"
	"myfavouritemovies/routes/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {
	router.POST("/users",handlers.AddUser)
    router.PATCH("/users/:id", handlers.UpdateUser)
	router.GET("/users/:id", handlers.ReadUser)
    router.POST("/users/:id/genres", handlers.AddFavoriteGenres)
	router.DELETE("/users/:id/genres", handlers.DeleteFavoriteGenre)
    router.POST("/users/:id/movies", handlers.AddFavoriteMovie)
	router.DELETE("/users/:id/movies", handlers.DeleteFavoriteMovie)
    router.PATCH("/users/:id/movies/toggle", handlers.ToggleWatchedStatus)
	router.GET("/users/:id/movies", apihandlers.GetMovies)
}
