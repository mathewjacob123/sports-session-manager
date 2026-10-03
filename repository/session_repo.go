package repository

import (
	"database/sql"
	"sports-session-manager/models"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(req models.CreateSessionRequest) (*models.Session, error) {
	var matchDate *time.Time
	var format *string
	var totalSlots *int

	if req.MatchDate != "" {
		t, err := time.Parse("2006-01-02", req.MatchDate)
		if err != nil {
			return nil, err
		}
		matchDate = &t
	}

	if req.Format != "" {
		format = &req.Format
		slots := 14
		if req.Format == "9v9" {
			slots = 18
		}
		totalSlots = &slots
	}

	status := "open"
	if matchDate == nil {
		status = "polling"
	}

	query := `
		INSERT INTO sessions (match_date, venue, format, total_slots, ground_fee, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, match_date, venue, format, total_slots, ground_fee, status, outside_notified, created_by, created_at
	`

	session := &models.Session{}
	err := r.db.QueryRow(query,
		matchDate, req.Venue, format, totalSlots, req.GroundFee, status, req.CreatedBy,
	).Scan(
		&session.ID, &session.MatchDate, &session.Venue, &session.Format,
		&session.TotalSlots, &session.GroundFee, &session.Status,
		&session.OutsideNotified, &session.CreatedBy, &session.CreatedAt,
	)

	return session, err
}

func (r *SessionRepository) GetSessionByID(id int) (*models.Session, error) {
	query := `
		SELECT id, match_date, venue, format, total_slots, ground_fee, status, outside_notified, created_by, created_at
		FROM sessions WHERE id = $1
	`
	session := &models.Session{}
	err := r.db.QueryRow(query, id).Scan(
		&session.ID, &session.MatchDate, &session.Venue, &session.Format,
		&session.TotalSlots, &session.GroundFee, &session.Status,
		&session.OutsideNotified, &session.CreatedBy, &session.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) GetSessions() ([]models.Session, error) {
	query := `
		SELECT id, match_date, venue, format, total_slots, ground_fee, status, outside_notified, created_by, created_at
		FROM sessions
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := []models.Session{}
	for rows.Next() {
		var s models.Session
		err := rows.Scan(
			&s.ID, &s.MatchDate, &s.Venue, &s.Format,
			&s.TotalSlots, &s.GroundFee, &s.Status,
			&s.OutsideNotified, &s.CreatedBy, &s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *SessionRepository) UpdateSessionStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE sessions SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (r *SessionRepository) MarkOutsideNotified(id int) error {
	_, err := r.db.Exec(`UPDATE sessions SET outside_notified = true WHERE id = $1`, id)
	return err
}

// Poll functions
func (r *SessionRepository) CreatePoll(sessionID int, options []models.PollOption) ([]models.DatePoll, error) {
	polls := []models.DatePoll{}
	for _, opt := range options {
		t, err := time.Parse("2006-01-02", opt.ProposedDate)
		if err != nil {
			return nil, err
		}
		var poll models.DatePoll
		err = r.db.QueryRow(`
			INSERT INTO date_polls (session_id, proposed_date, proposed_format)
			VALUES ($1, $2, $3)
			RETURNING id, session_id, proposed_date, proposed_format, created_at
		`, sessionID, t, opt.ProposedFormat).Scan(
			&poll.ID, &poll.SessionID, &poll.ProposedDate,
			&poll.ProposedFormat, &poll.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		polls = append(polls, poll)
	}
	return polls, nil
}

func (r *SessionRepository) GetPollsBySession(sessionID int) ([]models.DatePoll, error) {
	query := `
		SELECT dp.id, dp.session_id, dp.proposed_date, dp.proposed_format, dp.created_at,
		       COUNT(pv.id) as vote_count
		FROM date_polls dp
		LEFT JOIN poll_votes pv ON dp.id = pv.poll_id
		WHERE dp.session_id = $1
		GROUP BY dp.id
		ORDER BY vote_count DESC
	`
	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	polls := []models.DatePoll{}
	for rows.Next() {
		var p models.DatePoll
		err := rows.Scan(&p.ID, &p.SessionID, &p.ProposedDate, &p.ProposedFormat, &p.CreatedAt, &p.VoteCount)
		if err != nil {
			return nil, err
		}
		polls = append(polls, p)
	}
	return polls, nil
}

func (r *SessionRepository) VoteForPoll(pollID, playerID int) error {
	_, err := r.db.Exec(`
		INSERT INTO poll_votes (poll_id, player_id)
		VALUES ($1, $2)
	`, pollID, playerID)
	return err
}

func (r *SessionRepository) HasVoted(pollID, playerID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM poll_votes
		WHERE poll_id = $1 AND player_id = $2
	`, pollID, playerID).Scan(&count)
	return count > 0, err
}

func (r *SessionRepository) LockPollWinner(sessionID, pollID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// get winning poll details
	var matchDate time.Time
	var format string
	err = tx.QueryRow(`
		SELECT proposed_date, proposed_format FROM date_polls WHERE id = $1
	`, pollID).Scan(&matchDate, &format)
	if err != nil {
		return err
	}

	// calculate slots
	totalSlots := 14
	if format == "9v9" {
		totalSlots = 18
	}

	// update session with winning date and format
	_, err = tx.Exec(`
		UPDATE sessions
		SET match_date = $1, format = $2, total_slots = $3, status = 'open'
		WHERE id = $4
	`, matchDate, format, totalSlots, sessionID)
	if err != nil {
		return err
	}

	// delete losing poll options and their votes (cascade handles votes)
	_, err = tx.Exec(`
		DELETE FROM date_polls WHERE session_id = $1 AND id != $2
	`, sessionID, pollID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Slot functions
func (r *SessionRepository) GetSlotCounts(sessionID int) (*models.SlotCounts, error) {
	query := `
		SELECT
			COUNT(CASE WHEN slot_type = 'goalkeeper' AND status = 'confirmed' THEN 1 END) as gk_confirmed,
			COUNT(CASE WHEN slot_type = 'outfield'   AND status = 'confirmed' THEN 1 END) as out_confirmed,
			COUNT(CASE WHEN status = 'waitlisted' THEN 1 END) as waitlisted,
			COUNT(CASE WHEN status = 'confirmed' THEN 1 END) as total_confirmed
		FROM session_slots
		WHERE session_id = $1
	`
	counts := &models.SlotCounts{}
	err := r.db.QueryRow(query, sessionID).Scan(
		&counts.GoalkeeperConfirmed,
		&counts.OutfieldConfirmed,
		&counts.WaitlistCount,
		&counts.TotalConfirmed,
	)
	return counts, err
}

func (r *SessionRepository) GetSessionSlots(sessionID int) ([]models.SessionSlot, error) {
	query := `
		SELECT id, session_id, player_id, slot_type, player_type, status, confirmed_at, dropped_at
		FROM session_slots
		WHERE session_id = $1
		ORDER BY confirmed_at ASC
	`
	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slots := []models.SessionSlot{}
	for rows.Next() {
		var s models.SessionSlot
		err := rows.Scan(
			&s.ID, &s.SessionID, &s.PlayerID, &s.SlotType,
			&s.PlayerType, &s.Status, &s.ConfirmedAt, &s.DroppedAt,
		)
		if err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, nil
}

func (r *SessionRepository) PlayerAlreadyIn(sessionID, playerID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM session_slots
		WHERE session_id = $1 AND player_id = $2
		AND status IN ('confirmed', 'waitlisted')
	`, sessionID, playerID).Scan(&count)
	return count > 0, err
}

func (r *SessionRepository) JoinSession(sessionID, playerID int, playerType string, counts *models.SlotCounts, totalSlots int, position string) (*models.SessionSlot, error) {

	// determine slot type and whether they get in or waitlisted
	slotType := "outfield"
	status := "confirmed"

	if position == "goalkeeper" {
		slotType = "goalkeeper"
		if counts.GoalkeeperConfirmed >= 2 {
			// GK slots full — goes to waitlist as outfield
			slotType = "outfield"
			status = "waitlisted"
		}
	} else {
		// outfield player
		if counts.TotalConfirmed >= totalSlots {
			status = "waitlisted"
		}
	}

	slot := &models.SessionSlot{}
	err := r.db.QueryRow(`
		INSERT INTO session_slots (session_id, player_id, slot_type, player_type, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, session_id, player_id, slot_type, player_type, status, confirmed_at, dropped_at
	`, sessionID, playerID, slotType, playerType, status).Scan(
		&slot.ID, &slot.SessionID, &slot.PlayerID, &slot.SlotType,
		&slot.PlayerType, &slot.Status, &slot.ConfirmedAt, &slot.DroppedAt,
	)

	return slot, err
}

func (r *SessionRepository) DropFromSession(sessionID, playerID int, isMatchDay bool, groundFee float64, totalConfirmed int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// get player's current slot
	var slotID int
	var status, slotType string
	err = tx.QueryRow(`
		SELECT id, status, slot_type FROM session_slots
		WHERE session_id = $1 AND player_id = $2
		AND status IN ('confirmed', 'waitlisted')
	`, sessionID, playerID).Scan(&slotID, &status, &slotType)
	if err != nil {
		return err
	}

	// mark slot as dropped
	_, err = tx.Exec(`
		UPDATE session_slots
		SET status = 'dropped', dropped_at = NOW(), player_id = NULL
		WHERE id = $1
	`, slotID)
	if err != nil {
		return err
	}

	// if match day dropout and was confirmed (not waitlisted) — add penalty
	if isMatchDay && status == "confirmed" {
		perPersonFee := groundFee / float64(totalConfirmed)
		_, err = tx.Exec(`
			INSERT INTO penalties (player_id, session_id, amount_owed, reason)
			VALUES ($1, $2, $3, 'late_dropout')
			ON CONFLICT (player_id, session_id) DO NOTHING
		`, playerID, sessionID, perPersonFee)
		if err != nil {
			return err
		}

		// add strike for late dropout
		_, err = tx.Exec(`
			INSERT INTO strike_events (player_id, session_id, strike_type)
			VALUES ($1, $2, 'late_dropout')
		`, playerID, sessionID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE players SET strike_count = strike_count + 1 WHERE id = $1
		`, playerID)
		if err != nil {
			return err
		}
	}

	// promote first waitlisted player if confirmed slot was freed
	if status == "confirmed" {
		_, err = tx.Exec(`
			UPDATE session_slots
			SET status = 'confirmed', confirmed_at = NOW()
			WHERE id = (
				SELECT id FROM session_slots
				WHERE session_id = $1 AND status = 'waitlisted'
				ORDER BY created_at ASC
				LIMIT 1
			)
		`, sessionID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}