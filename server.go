package server

import (
	"myfavouritemovies/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func CreateServer() *gin.Engine {
	router := gin.Default()
	router.Use(sessions.Sessions("my_session", cookie.NewStore([]byte("secret"))))
	routes.SetUpRoutes(router)
	return router
}
