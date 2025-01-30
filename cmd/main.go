package main

import (
	"fmt"
	"log"
	server "myfavouritemovies"
	"myfavouritemovies/database"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err:=godotenv.Load()
	if err != nil {
		log.Println("Not found .env file")
	}
	port:=os.Getenv("SRVR_PORT")
	db := database.InitDB()
	if db != nil {
		fmt.Println("Database tables created successfully!")
	} else {
		log.Fatal("Failed to initialize the database.")
	}

	server := server.CreateServer()

	if err := server.Run(":"+port); err != nil {
		log.Fatal("Server run failed: ", err)
	} else {
		fmt.Println("SERVER RUNNING")
	}
}

