CREATE TABLE mistakes (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL,
    question_id    UUID NOT NULL,
    language_id    UUID NOT NULL,
    user_answer    TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    reviewed       BOOLEAN NOT NULL DEFAULT false,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_mistakes_user ON mistakes(user_id, created_at DESC);
CREATE INDEX idx_mistakes_language ON mistakes(user_id, language_id);
CREATE INDEX idx_mistakes_unreviewed ON mistakes(user_id, reviewed) WHERE reviewed = false;
