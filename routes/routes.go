package routes

import (
	"net/http"

	"github.com/gaurangi/mvc_assignment/controllers"
	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(queries *db.Queries, dbPool *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/register", controllers.RegisterPlayer(queries))
	mux.HandleFunc("/api/login", controllers.LoginPlayer(queries))
	mux.HandleFunc("/api/player/dashboard", middleware.RequireAuth(controllers.GetDashboard(queries)))

	mux.HandleFunc("/api/village", middleware.RequireAuth(controllers.GetVillage(queries)))
	mux.HandleFunc("/api/village/buildings", middleware.RequireAuth(controllers.PurchaseBuilding(queries)))
	mux.HandleFunc("/api/village/buildings/upgrade", middleware.RequireAuth(controllers.UpgradeBuilding(queries)))
	mux.HandleFunc("/api/village/buildings/upgrade/cost", middleware.RequireAuth(controllers.GetBuildingUpgradeCost(queries)))
	mux.HandleFunc("/api/village/buildings/complete", middleware.RequireAuth(controllers.CompleteBuild(queries)))
	mux.HandleFunc("/api/village/move-building", middleware.RequireAuth(controllers.MoveBuilding(dbPool)))
	mux.HandleFunc("/api/village/upgrade", middleware.RequireAuth(controllers.UpgradePlayerVillage(dbPool)))
	mux.HandleFunc("/api/village/upgrade/cost", middleware.RequireAuth(controllers.GetVillageUpgradeCost(dbPool)))
	mux.HandleFunc("/api/village/collect", middleware.RequireAuth(controllers.CollectResources(queries)))

	mux.HandleFunc("/api/shop/catalog", middleware.RequireAuth(controllers.GetShopCatalog(queries)))

	mux.HandleFunc("/api/army/catalog", middleware.RequireAuth(controllers.GetArmyCatalog(queries)))
	mux.HandleFunc("/api/army/train", middleware.RequireAuth(controllers.TrainTroops(dbPool)))
	mux.HandleFunc("/api/army/status", middleware.RequireAuth(controllers.GetArmyStatus(queries)))

	mux.HandleFunc("/api/battle/match", middleware.RequireAuth(controllers.GetMatch(queries)))
	mux.HandleFunc("/api/battle/attack", middleware.RequireAuth(controllers.AttackOpponent(dbPool)))
	mux.HandleFunc("/api/battle/replay", middleware.RequireAuth(controllers.GetBattleReplay(queries)))

	return mux
}
