package handlers

import (
	"net/http"
	"strconv"
	"sports-session-manager/services"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(teamService *services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

func (h *TeamHandler) GenerateTeams(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	teams, err := h.teamService.GenerateTeams(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}