CREATE TABLE enrollments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    language_id UUID NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, language_id)
);

CREATE INDEX idx_enrollments_user ON enrollments(user_id) WHERE is_active = true;
CREATE INDEX idx_enrollments_language ON enrollments(language_id);
