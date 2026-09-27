CREATE TABLE IF NOT EXISTS penalties (
    id          SERIAL PRIMARY KEY,
    player_id   INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    session_id  INTEGER NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    amount_owed NUMERIC(10,2) NOT NULL,
    reason      VARCHAR(50) NOT NULL DEFAULT 'late_dropout',
    is_paid     BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(player_id, session_id)
);