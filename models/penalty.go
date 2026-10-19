package models

import "time"

type Penalty struct {
	ID          int       `json:"id"`
	PlayerID    int       `json:"player_id"`
	SessionID   int       `json:"session_id"`
	AmountOwed  float64   `json:"amount_owed"`
	Reason      string    `json:"reason"`
	IsPaid      bool      `json:"is_paid"`
	CreatedAt   time.Time `json:"created_at"`
}

type PenaltyWithDetails struct {
	Penalty
	PlayerName  string  `json:"player_name"`
	Venue       string  `json:"venue"`
	TotalOwed   float64 `json:"total_owed,omitempty"`
}