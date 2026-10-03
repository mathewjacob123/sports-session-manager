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
	playerRepo  := repository.NewPlayerRepository(database)
	sessionRepo := repository.NewSessionRepository(database)

	// services
	playerService  := services.NewPlayerService(playerRepo)
	sessionService := services.NewSessionService(sessionRepo, playerRepo)

	// handlers
	playerHandler  := handlers.NewPlayerHandler(playerService)
	sessionHandler := handlers.NewSessionHandler(sessionService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// player routes
	players := r.Group("/players")
	{
		players.POST("",               playerHandler.CreatePlayer)
		players.GET("",                playerHandler.GetPlayers)
		players.GET("/:id",            playerHandler.GetPlayerByID)
		players.PUT("/:id",            playerHandler.UpdatePlayer)
		players.POST("/:id/boot",      playerHandler.BootPlayer)
		players.POST("/:id/strike",    playerHandler.AddStrike)
		players.GET("/:id/strikes",    playerHandler.GetPlayerStrikes)
	}

	// session routes
	sessions := r.Group("/sessions")
	{
		sessions.POST("",                              sessionHandler.CreateSession)
		sessions.GET("",                               sessionHandler.GetSessions)
		sessions.GET("/:id",                           sessionHandler.GetSessionByID)
		sessions.POST("/:id/confirm",                  sessionHandler.ConfirmSession)
		sessions.POST("/:id/complete",                 sessionHandler.CompleteSession)
		sessions.POST("/:id/join",                     sessionHandler.JoinSession)
		sessions.POST("/:id/drop",                     sessionHandler.DropFromSession)
		sessions.GET("/:id/check-outside",             sessionHandler.CheckOutsideNotification)
		sessions.POST("/:id/polls",                    sessionHandler.CreatePoll)
		sessions.GET("/:id/polls",                     sessionHandler.GetPolls)
		sessions.POST("/:id/polls/:pollId/vote",       sessionHandler.VoteForPoll)
		sessions.POST("/:id/polls/:pollId/lock",       sessionHandler.LockPollWinner)
	}

	log.Println("Server running on port", cfg.Port)
	r.Run(":" + cfg.Port)
}