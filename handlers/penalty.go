package handlers

import (
	"net/http"
	"strconv"
	"sports-session-manager/services"

	"github.com/gin-gonic/gin"
)

type PenaltyHandler struct {
	penaltyService *services.PenaltyService
}

func NewPenaltyHandler(penaltyService *services.PenaltyService) *PenaltyHandler {
	return &PenaltyHandler{penaltyService: penaltyService}
}

func (h *PenaltyHandler) GetPenaltiesByPlayer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}

	penalties, totalOwed, err := h.penaltyService.GetPenaltiesByPlayer(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"penalties":     penalties,
		"total_owed":    totalOwed,
		"penalty_count": len(penalties),
	})
}

func (h *PenaltyHandler) GetPenaltiesBySession(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	penalties, totalOwed, err := h.penaltyService.GetPenaltiesBySession(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"penalties":     penalties,
		"total_owed":    totalOwed,
		"penalty_count": len(penalties),
	})
}

func (h *PenaltyHandler) MarkAsPaid(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid penalty id"})
		return
	}

	err = h.penaltyService.MarkAsPaid(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "penalty marked as paid"})
}