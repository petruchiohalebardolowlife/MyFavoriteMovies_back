package main

import (
	"fmt"
	"log"
	"myfavouritemovies/database"
	"myfavouritemovies/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.InitDB()

	if db != nil {
		fmt.Println("Database tables created successfully!")
	} else {
		log.Fatal("Failed to initialize the database.")
	}

	router:=gin.Default()

	routes.SetUpRoutes(router)

	if err := router.Run(":8081"); err != nil {
        log.Fatal("Server run failed: ", err)
    } else {
		fmt.Println("SERVER RUNNING")
	}
}
