CREATE TABLE IF NOT EXISTS players (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    phone       VARCHAR(20) UNIQUE NOT NULL,
    position    VARCHAR(50) NOT NULL CHECK (position IN ('goalkeeper', 'defender', 'midfielder', 'forward')),
    skill_rating INTEGER NOT NULL DEFAULT 5 CHECK (skill_rating BETWEEN 1 AND 10),
    player_type VARCHAR(20) NOT NULL DEFAULT 'group' CHECK (player_type IN ('group', 'outside')),
    strike_count INTEGER NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);