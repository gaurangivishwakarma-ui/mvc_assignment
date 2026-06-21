package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gaurangi/mvc_assignment/config"
	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/routes"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Relying on system environment variables.")
	}

	dbPool := config.InitDB()
	defer dbPool.Close()

	queries := db.New(dbPool)

	mux := routes.SetupRoutes(queries, dbPool)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started successfully on port %s\n", port)

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server crashed: %v\n", err)
	}
}
