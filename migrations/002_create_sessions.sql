CREATE TABLE IF NOT EXISTS sessions (
    id              SERIAL PRIMARY KEY,
    match_date      DATE NOT NULL,
    venue           VARCHAR(255),
    format          VARCHAR(10) NOT NULL DEFAULT '7v7' CHECK (format IN ('7v7', '9v9')),
    total_slots     INTEGER NOT NULL DEFAULT 14,
    ground_fee      NUMERIC(10,2) NOT NULL DEFAULT 0,
    status          VARCHAR(20) NOT NULL DEFAULT 'polling' 
                    CHECK (status IN ('polling', 'confirmed', 'completed', 'cancelled')),
    created_by      INTEGER NOT NULL REFERENCES players(id),
    outside_notified BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);