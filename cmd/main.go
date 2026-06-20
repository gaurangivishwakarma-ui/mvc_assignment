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

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	queries := db.New(dbPool)
	_ = queries

	http.HandleFunc("/api/register", controllers.RegisterPlayer(queries))
	http.HandleFunc("/api/login", controllers.LoginPlayer(queries))
	http.HandleFunc("/api/player/dashboard", middleware.RequireAuth(controllers.GetDashboard(queries)))

	http.HandleFunc("/api/village", middleware.RequireAuth(controllers.GetVillage(queries)))
	http.HandleFunc("/api/village/buildings", middleware.RequireAuth(controllers.PurchaseBuilding(queries)))
	http.HandleFunc("/api/village/buildings/upgrade", middleware.RequireAuth(controllers.UpgradeBuilding(queries)))
	http.HandleFunc("/api/village/move-building", middleware.RequireAuth(controllers.MoveBuilding(dbPool)))
	http.HandleFunc("/api/village/upgrade", middleware.RequireAuth(controllers.UpgradePlayerVillage(dbPool)))
	http.HandleFunc("/api/village/collect", middleware.RequireAuth(controllers.CollectResources(queries)))

	http.HandleFunc("/api/army/catalog", middleware.RequireAuth(controllers.GetArmyCatalog(queries)))
	http.HandleFunc("/api/army/train", middleware.RequireAuth(controllers.TrainTroops(dbPool)))
	http.HandleFunc("/api/army/status", middleware.RequireAuth(controllers.GetArmyStatus(queries)))

	http.HandleFunc("/api/battle/match", middleware.RequireAuth(controllers.GetMatch(queries)))
	http.HandleFunc("/api/battle/attack", middleware.RequireAuth(controllers.AttackOpponent(dbPool)))
	http.HandleFunc("/api/battle/replay", middleware.RequireAuth(controllers.GetBattleReplay(queries)))

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started successfully on port %s\n", port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server crashed: %v\n", err)
	}
}
