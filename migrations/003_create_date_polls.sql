CREATE TABLE IF NOT EXISTS date_polls (
    id          SERIAL PRIMARY KEY,
    session_id  INTEGER NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    proposed_date DATE NOT NULL,
    vote_count  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS poll_votes (
    id          SERIAL PRIMARY KEY,
    poll_id     INTEGER NOT NULL REFERENCES date_polls(id) ON DELETE CASCADE,
    player_id   INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(poll_id, player_id)
);