CREATE TABLE IF NOT EXISTS session_slots (
    id          SERIAL PRIMARY KEY,
    session_id  INTEGER NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    player_id   INTEGER REFERENCES players(id) ON DELETE SET NULL,
    slot_type   VARCHAR(20) NOT NULL DEFAULT 'outfield' 
                CHECK (slot_type IN ('goalkeeper', 'outfield')),
    player_type VARCHAR(20) CHECK (player_type IN ('group', 'outside')),
    status      VARCHAR(20) NOT NULL DEFAULT 'open' 
                CHECK (status IN ('open', 'confirmed', 'dropped', 'waitlisted')),
    confirmed_at TIMESTAMP,
    dropped_at  TIMESTAMP,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);