package routes

import (
	"myfavouritemovies/routes/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {
	router.POST("/users",handlers.AddUser)
    router.PATCH("/users/:id", handlers.UpdateUser)
	router.GET("/users/:id", handlers.ReadUser)
    router.POST("/users/:id/genres", handlers.AddOrDeleteFavoriteGenre)
    router.POST("/users/:id/movies", handlers.AddFavoriteMovie)
	router.DELETE("/users/:id/movies", handlers.DeleteFavoriteMovie)
    router.PATCH("/users/:id/movies/toggle", handlers.ToggleWatchedStatus)
    router.POST("/genres", handlers.AddGenres)
}
