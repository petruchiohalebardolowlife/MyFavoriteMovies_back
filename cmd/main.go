package main

import (
	"fmt"
	"log"
	server "myfavouritemovies"
	config "myfavouritemovies/configs"
	"myfavouritemovies/database"
)

func main() {
	config.LoadConfig()

	db := database.InitDB()
	if db != nil {
		fmt.Println("Database tables created successfully!")
	} else {
		log.Fatal("Failed to initialize the database.")
	}

	server := server.CreateServer()

	if err := server.Run(":"+config.SRVR_PORT); err != nil {
		log.Fatal("Server run failed: ", err)
	} else {
		fmt.Println("SERVER RUNNING")
	}
}

