package models

import "time"

type Player struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Position    string    `json:"position"`
	SkillRating int       `json:"skill_rating"`
	PlayerType  string    `json:"player_type"`
	StrikeCount int       `json:"strike_count"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreatePlayerRequest struct {
	Name        string `json:"name"         binding:"required"`
	Phone       string `json:"phone"        binding:"required"`
	Position    string `json:"position"     binding:"required,oneof=goalkeeper defender midfielder forward"`
	SkillRating int    `json:"skill_rating" binding:"required,min=1,max=10"`
	PlayerType  string `json:"player_type"  binding:"required,oneof=group outside"`
}

type UpdatePlayerRequest struct {
	Name        string `json:"name"`
	Position    string `json:"position"     binding:"omitempty,oneof=goalkeeper defender midfielder forward"`
	SkillRating int    `json:"skill_rating" binding:"omitempty,min=1,max=10"`
}