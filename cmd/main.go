package main

import (
	"fmt"
	"log"
	server "myfavouritemovies"
	config "myfavouritemovies/configs"
	"myfavouritemovies/database"
	"myfavouritemovies/repository"
	api "myfavouritemovies/routes/apihandlers"
)

func main() {
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
		fmt.Println("Server running")
	}

  genres, err := api.FetchGenres()
  if err != nil {
    log.Fatalf("Failed to fetch genres: %v", err)
  }
  if err := repository.SaveGenresToDB(db, genres); err != nil {
		log.Fatalf("Failed to add genres: %v", err)
	}
}

