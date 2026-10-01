package repository

import (
	"database/sql"
	"sports-session-manager/models"
)

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) CreatePlayer(req models.CreatePlayerRequest) (*models.Player, error) {
	query := `
		INSERT INTO players (name, phone, position, skill_rating, player_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, phone, position, skill_rating, player_type, strike_count, is_active, created_at
	`
	player := &models.Player{}
	err := r.db.QueryRow(query,
		req.Name, req.Phone, req.Position, req.SkillRating, req.PlayerType,
	).Scan(
		&player.ID, &player.Name, &player.Phone, &player.Position,
		&player.SkillRating, &player.PlayerType, &player.StrikeCount,
		&player.IsActive, &player.CreatedAt,
	)
	return player, err
}

func (r *PlayerRepository) GetPlayers(activeOnly bool) ([]models.Player, error) {
	query := `
		SELECT id, name, phone, position, skill_rating, player_type, strike_count, is_active, created_at
		FROM players
	`
	if activeOnly {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY name ASC"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := []models.Player{}
	for rows.Next() {
		var p models.Player
		err := rows.Scan(
			&p.ID, &p.Name, &p.Phone, &p.Position,
			&p.SkillRating, &p.PlayerType, &p.StrikeCount,
			&p.IsActive, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}

func (r *PlayerRepository) GetPlayerByID(id int) (*models.Player, error) {
	query := `
		SELECT id, name, phone, position, skill_rating, player_type, strike_count, is_active, created_at
		FROM players WHERE id = $1
	`
	player := &models.Player{}
	err := r.db.QueryRow(query, id).Scan(
		&player.ID, &player.Name, &player.Phone, &player.Position,
		&player.SkillRating, &player.PlayerType, &player.StrikeCount,
		&player.IsActive, &player.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return player, nil
}

func (r *PlayerRepository) GetPlayerByPhone(phone string) (*models.Player, error) {
	query := `
		SELECT id, name, phone, position, skill_rating, player_type, strike_count, is_active, created_at
		FROM players WHERE phone = $1
	`
	player := &models.Player{}
	err := r.db.QueryRow(query, phone).Scan(
		&player.ID, &player.Name, &player.Phone, &player.Position,
		&player.SkillRating, &player.PlayerType, &player.StrikeCount,
		&player.IsActive, &player.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return player, nil
}

func (r *PlayerRepository) UpdatePlayer(id int, req models.UpdatePlayerRequest) (*models.Player, error) {
	query := `
		UPDATE players
		SET name = COALESCE(NULLIF($1, ''), name),
		    position = COALESCE(NULLIF($2, ''), position),
		    skill_rating = CASE WHEN $3 = 0 THEN skill_rating ELSE $3 END
		WHERE id = $4
		RETURNING id, name, phone, position, skill_rating, player_type, strike_count, is_active, created_at
	`
	player := &models.Player{}
	err := r.db.QueryRow(query,
		req.Name, req.Position, req.SkillRating, id,
	).Scan(
		&player.ID, &player.Name, &player.Phone, &player.Position,
		&player.SkillRating, &player.PlayerType, &player.StrikeCount,
		&player.IsActive, &player.CreatedAt,
	)
	return player, err
}

func (r *PlayerRepository) BootPlayer(id int) error {
	query := `UPDATE players SET is_active = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PlayerRepository) AddStrike(playerID, sessionID int, strikeType string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// insert strike event
	_, err = tx.Exec(`
		INSERT INTO strike_events (player_id, session_id, strike_type)
		VALUES ($1, $2, $3)
	`, playerID, sessionID, strikeType)
	if err != nil {
		return err
	}

	// increment strike count on player
	var newCount int
	err = tx.QueryRow(`
		UPDATE players SET strike_count = strike_count + 1
		WHERE id = $1
		RETURNING strike_count
	`, playerID).Scan(&newCount)
	if err != nil {
		return err
	}

	// auto boot at 3 strikes
	if newCount >= 3 {
		_, err = tx.Exec(`
			UPDATE players SET is_active = false WHERE id = $1
		`, playerID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PlayerRepository) GetPlayerStrikes(playerID int) ([]map[string]interface{}, error) {
	query := `
		SELECT se.id, se.strike_type, se.created_at, s.match_date
		FROM strike_events se
		JOIN sessions s ON se.session_id = s.id
		WHERE se.player_id = $1
		ORDER BY se.created_at DESC
	`
	rows, err := r.db.Query(query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	strikes := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var strikeType, createdAt, matchDate string
		rows.Scan(&id, &strikeType, &createdAt, &matchDate)
		strikes = append(strikes, map[string]interface{}{
			"id":          id,
			"strike_type": strikeType,
			"created_at":  createdAt,
			"match_date":  matchDate,
		})
	}
	return strikes, nil
}