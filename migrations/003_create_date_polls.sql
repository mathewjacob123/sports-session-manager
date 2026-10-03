CREATE TABLE IF NOT EXISTS date_polls (
    id              SERIAL PRIMARY KEY,
    session_id      INTEGER NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    proposed_date   DATE NOT NULL,
    proposed_format VARCHAR(10) NOT NULL CHECK (proposed_format IN ('7v7', '9v9')),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);