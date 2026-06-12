package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/gaurangi/mvc_assignment/controllers"
	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Relying on system environment variables.")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	// 3. Connect to PostgreSQL
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("❌ Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close() // Ensures the connection closes when the server shuts down
	log.Println("✅ Successfully connected to the PostgreSQL database!")

	// 4. Initialize your generated sqlc queries
	// Now your sqlc functions have a live database connection to use!
	queries := db.New(dbPool)
	_ = queries
	// 5. Setup Routes (Connecting the Waiters to the Front Door)
	// Example: http.HandleFunc("/api/register", controllers.RegisterHandler(queries))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Game server is running cleanly!"))
	})

	http.HandleFunc("/api/register", controllers.RegisterPlayer(queries))
	http.HandleFunc("/api/login", controllers.LoginPlayer(queries))

	http.HandleFunc("/api/village", middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		// Because the Bouncer let them through, we can safely pull their ID directly from the context!
		playerID := r.Context().Value(middleware.PlayerIDKey).(string)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"message": "Welcome to your village! Your secure ID is %s"}`, playerID)))
	}))

	// 6. Start the Server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Starting game server on port %s...\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("❌ Server crashed: %v\n", err)
	}
}
