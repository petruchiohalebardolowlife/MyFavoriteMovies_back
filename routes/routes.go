package routes

import (
	"myfavouritemovies/routes/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {
	router.POST("/movies",handlers.AddMovie)
	router.POST("/users",handlers.AddUser)
    router.PATCH("/users/:id", handlers.UpdateUser)
	router.GET("/users/:id", handlers.ReadUser)
    router.POST("/users/:id/genres", handlers.AddFavouriteGenre)
    router.POST("/users/:id/movies", handlers.AddFavouriteMovie)
    router.PATCH("/users/:id/movies/toggle", handlers.ToggleWatchedStatus)
    router.POST("/genres", handlers.AddGenre)
}
