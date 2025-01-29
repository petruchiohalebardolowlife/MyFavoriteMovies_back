package main

import (
	"fmt"
	"log"
	server "myfavouritemovies"
	"myfavouritemovies/database"
)

func main() {
	db := database.InitDB()
	if db != nil {
		fmt.Println("Database tables created successfully!")
	} else {
		log.Fatal("Failed to initialize the database.")
	}

	server := server.CreateServer()

	if err := server.Run(":8081"); err != nil {
		log.Fatal("Server run failed: ", err)
	} else {
		fmt.Println("SERVER RUNNING")
	}
}

