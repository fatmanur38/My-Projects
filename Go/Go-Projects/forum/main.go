package main

import (
	
	"log"
	"net/http"

	
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver

	"forum/handlers"
	"forum/utils"
)



func main() {
	utils.InitDatabase()   // Initialize the global database connection
	defer utils.Db.Close() // Ensure the database connection is closed when main exits
	utils.CreateTables()   // Create database tables

	handlers.SetupRoutes() // Set up HTTP server

	err := http.ListenAndServe(":8040", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

