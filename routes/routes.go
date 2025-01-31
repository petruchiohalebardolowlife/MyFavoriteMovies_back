package routes

import (
	"encoding/gob"
	"myfavouritemovies/routes/handlers"
	"myfavouritemovies/structs"

	"github.com/gin-gonic/gin"
)


func SetUpRoutes(router *gin.Engine) {
	gob.Register(structs.User{})  

	router.POST("/login",handlers.Login)
	router.POST("/users",handlers.AddUser)
    router.PATCH("/users/:id", handlers.UpdateUser)
    router.POST("/favgenres", handlers.AddFavoriteGenre)
	router.DELETE("/favgenres", handlers.DeleteFavoriteGenre)
    router.POST("/movies", handlers.AddFavoriteMovie)
	router.DELETE("/movies", handlers.DeleteFavoriteMovie)
    router.PATCH("/movies/toggle", handlers.ToggleWatchedStatus)
    router.POST("/genres", handlers.AddGenres)
}
