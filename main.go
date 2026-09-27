package main

import (
	"sports-session-manager/config"
	"sports-session-manager/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database := db.Connect(cfg)
	defer database.Close()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Server running on port", cfg.Port)
	r.Run(":" + cfg.Port)
}