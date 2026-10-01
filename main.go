package main

import (
	"sports-session-manager/config"
	"sports-session-manager/db"
	"sports-session-manager/handlers"
	"sports-session-manager/repository"
	"sports-session-manager/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database := db.Connect(cfg)
	defer database.Close()

	// repositories
	playerRepo := repository.NewPlayerRepository(database)

	// services
	playerService := services.NewPlayerService(playerRepo)

	// handlers
	playerHandler := handlers.NewPlayerHandler(playerService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// player routes
	players := r.Group("/players")
	{
		players.POST("",              playerHandler.CreatePlayer)
		players.GET("",               playerHandler.GetPlayers)
		players.GET("/:id",           playerHandler.GetPlayerByID)
		players.PUT("/:id",           playerHandler.UpdatePlayer)
		players.POST("/:id/boot",     playerHandler.BootPlayer)
		players.POST("/:id/strike",   playerHandler.AddStrike)
		players.GET("/:id/strikes",   playerHandler.GetPlayerStrikes)
	}

	log.Println("Server running on port", cfg.Port)
	r.Run(":" + cfg.Port)
}