CREATE TABLE units (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    language_id UUID NOT NULL,
    title       TEXT NOT NULL,
    description TEXT DEFAULT '',
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_units_language ON units(language_id, sort_order);
