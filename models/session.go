package models

import "time"

type Session struct {
	ID               int        `json:"id"`
	MatchDate        *time.Time `json:"match_date"`
	Venue            string     `json:"venue"`
	Format           *string    `json:"format"`
	TotalSlots       *int       `json:"total_slots"`
	GroundFee        float64    `json:"ground_fee"`
	Status           string     `json:"status"`
	OutsideNotified  bool       `json:"outside_notified"`
	CreatedBy        int        `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
}

type CreateSessionRequest struct {
	MatchDate string  `json:"match_date"`
	Venue     string  `json:"venue"      binding:"required"`
	Format    string  `json:"format"     binding:"omitempty,oneof=7v7 9v9"`
	GroundFee float64 `json:"ground_fee" binding:"required,min=0"`
	CreatedBy int     `json:"created_by" binding:"required"`
}

type CreatePollRequest struct {
	Options []PollOption `json:"options" binding:"required,min=2"`
}

type PollOption struct {
	ProposedDate   string `json:"proposed_date"   binding:"required"`
	ProposedFormat string `json:"proposed_format" binding:"required,oneof=7v7 9v9"`
}

type DatePoll struct {
	ID             int       `json:"id"`
	SessionID      int       `json:"session_id"`
	ProposedDate   time.Time `json:"proposed_date"`
	ProposedFormat string    `json:"proposed_format"`
	VoteCount      int       `json:"vote_count"`
	CreatedAt      time.Time `json:"created_at"`
}

type SessionSlot struct {
	ID          int        `json:"id"`
	SessionID   int        `json:"session_id"`
	PlayerID    *int       `json:"player_id"`
	SlotType    string     `json:"slot_type"`
	PlayerType  string     `json:"player_type"`
	Status      string     `json:"status"`
	ConfirmedAt time.Time  `json:"confirmed_at"`
	DroppedAt   *time.Time `json:"dropped_at"`
}

type SessionResponse struct {
	Session
	ConfirmedCount int           `json:"confirmed_count"`
	WaitlistCount  int           `json:"waitlist_count"`
	Slots          []SessionSlot `json:"slots,omitempty"`
}

type SlotCounts struct {
	GoalkeeperConfirmed int
	OutfieldConfirmed   int
	WaitlistCount       int
	TotalConfirmed      int
}