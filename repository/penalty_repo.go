package repository

import (
	"database/sql"
	"sports-session-manager/models"
)

type PenaltyRepository struct {
	db *sql.DB
}

func NewPenaltyRepository(db *sql.DB) *PenaltyRepository {
	return &PenaltyRepository{db: db}
}

func (r *PenaltyRepository) GetPenaltiesByPlayer(playerID int) ([]models.PenaltyWithDetails, error) {
	query := `
		SELECT p.id, p.player_id, p.session_id, p.amount_owed, p.reason, p.is_paid, p.created_at,
		       pl.name as player_name, s.venue
		FROM penalties p
		JOIN players pl ON p.player_id = pl.id
		JOIN sessions s  ON p.session_id = s.id
		WHERE p.player_id = $1
		ORDER BY p.created_at DESC
	`
	return r.scanPenalties(query, playerID)
}

func (r *PenaltyRepository) GetPenaltiesBySession(sessionID int) ([]models.PenaltyWithDetails, error) {
	query := `
		SELECT p.id, p.player_id, p.session_id, p.amount_owed, p.reason, p.is_paid, p.created_at,
		       pl.name as player_name, s.venue
		FROM penalties p
		JOIN players pl ON p.player_id = pl.id
		JOIN sessions s  ON p.session_id = s.id
		WHERE p.session_id = $1
		ORDER BY p.created_at DESC
	`
	return r.scanPenalties(query, sessionID)
}

func (r *PenaltyRepository) GetPenaltyByID(id int) (*models.Penalty, error) {
	query := `
		SELECT id, player_id, session_id, amount_owed, reason, is_paid, created_at
		FROM penalties WHERE id = $1
	`
	p := &models.Penalty{}
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.PlayerID, &p.SessionID,
		&p.AmountOwed, &p.Reason, &p.IsPaid, &p.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (r *PenaltyRepository) MarkAsPaid(id int) error {
	_, err := r.db.Exec(`
		UPDATE penalties SET is_paid = true WHERE id = $1
	`, id)
	return err
}

func (r *PenaltyRepository) scanPenalties(query string, arg int) ([]models.PenaltyWithDetails, error) {
	rows, err := r.db.Query(query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	penalties := []models.PenaltyWithDetails{}
	for rows.Next() {
		var p models.PenaltyWithDetails
		err := rows.Scan(
			&p.ID, &p.PlayerID, &p.SessionID,
			&p.AmountOwed, &p.Reason, &p.IsPaid, &p.CreatedAt,
			&p.PlayerName, &p.Venue,
		)
		if err != nil {
			return nil, err
		}
		penalties = append(penalties, p)
	}
	return penalties, nil
}