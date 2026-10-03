CREATE TABLE IF NOT EXISTS sessions (
    id                  SERIAL PRIMARY KEY,
    match_date          DATE,
    venue               VARCHAR(255),
    format              VARCHAR(10) CHECK (format IN ('7v7', '9v9')),
    total_slots         INTEGER,
    ground_fee          NUMERIC(10,2) NOT NULL DEFAULT 0,
    status              VARCHAR(20) NOT NULL DEFAULT 'open'
                        CHECK (status IN ('open', 'polling', 'confirmed', 'completed', 'cancelled')),
    outside_notified    BOOLEAN NOT NULL DEFAULT false,
    created_by          INTEGER NOT NULL REFERENCES players(id),
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);