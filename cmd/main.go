package main

import (
	"fmt"
	"log"
	server "myfavouritemovies"
	config "myfavouritemovies/configs"
	"myfavouritemovies/database"
	"myfavouritemovies/repository"
	"myfavouritemovies/structs"
)

func main() {
	config.LoadConfig()

	db := database.InitDB()
	if db != nil {
		fmt.Println("Database tables created successfully!")
	} else {
		log.Fatal("Failed to initialize the database.")
	}
  genres := []structs.Genre{
    {ID: 12, Name: "Adventure"},
    {ID: 14, Name: "Fantasy"},
    {ID: 16, Name: "Animation"},
    {ID: 18, Name: "Drama"},
}
if err := repository.AddGenres(genres); err != nil {
  log.Fatalf("Failed to add genres: %v", err)
}
	server := server.CreateServer()

	if err := server.Run(":"+config.SRVR_PORT); err != nil {
		log.Fatal("Server run failed: ", err)
	} else {
		fmt.Println("SERVER RUNNING")
	}
}

