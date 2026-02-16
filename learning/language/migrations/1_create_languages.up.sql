CREATE TABLE languages (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,       -- "Yoruba", "Igbo", "Hausa"
    code        TEXT NOT NULL UNIQUE,       -- "yo", "ig", "ha"
    description TEXT DEFAULT '',
    flag_emoji  TEXT DEFAULT '',
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_languages_code ON languages(code);
CREATE INDEX idx_languages_active ON languages(is_active) WHERE is_active = true;
