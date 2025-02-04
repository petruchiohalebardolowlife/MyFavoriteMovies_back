package main

import (
	"fmt"
	"log"
	server "myfavouritemovies"
	config "myfavouritemovies/configs"
	"myfavouritemovies/database"
	"myfavouritemovies/repository"
	tmdb "myfavouritemovies/service"
	"myfavouritemovies/structs"
)

var movieFilter = structs.MovieFilter{
  GenreIDs: []int{18, 28},
  Rating:   "6", 
  Year:     2020,        
  Page:     1,      
}

var sonic_id = 454626

func main() {
	db := database.InitDB()
	if db != nil {
		fmt.Println("Database tables created successfully!")
	} else {
		log.Fatal("Failed to initialize the database.")
	}

  genres, err := tmdb.FetchGenres()
  if err != nil {
    log.Fatalf("Failed to fetch genres: %v", err)
  }
  if err := repository.SaveGenresToDB(db, genres); err != nil {
		log.Fatalf("Failed to add genres: %v", err)
	}

	server := server.CreateServer()
  if err := server.Run(":"+config.SRVR_PORT); err != nil {
		log.Fatal("Server run failed: ", err)
	} else {
		fmt.Println("Server running")
	}
}


