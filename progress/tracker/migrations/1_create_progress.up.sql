CREATE TABLE lesson_progress (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL,
    lesson_id    UUID NOT NULL,
    status       TEXT NOT NULL DEFAULT 'not_started',  -- "not_started" | "in_progress" | "completed"
    xp_earned    INT NOT NULL DEFAULT 0,
    attempts     INT NOT NULL DEFAULT 0,
    best_score   FLOAT NOT NULL DEFAULT 0,
    completed_at TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, lesson_id)
);

CREATE INDEX idx_progress_user ON lesson_progress(user_id);
CREATE INDEX idx_progress_lesson ON lesson_progress(lesson_id);
CREATE INDEX idx_progress_status ON lesson_progress(status);
