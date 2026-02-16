CREATE TABLE lesson_sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    lesson_id       UUID NOT NULL,
    status          TEXT NOT NULL DEFAULT 'in_progress',  -- "in_progress" | "completed" | "abandoned"
    current_index   INT NOT NULL DEFAULT 0,               -- current question index
    correct_count   INT NOT NULL DEFAULT 0,
    incorrect_count INT NOT NULL DEFAULT 0,
    xp_earned       INT NOT NULL DEFAULT 0,
    started_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user ON lesson_sessions(user_id);
CREATE INDEX idx_sessions_lesson ON lesson_sessions(lesson_id);
CREATE INDEX idx_sessions_status ON lesson_sessions(status) WHERE status = 'in_progress';

CREATE TABLE session_answers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id      UUID NOT NULL REFERENCES lesson_sessions(id) ON DELETE CASCADE,
    question_id     UUID NOT NULL,
    user_answer     TEXT NOT NULL,
    is_correct      BOOLEAN NOT NULL,
    similarity      FLOAT DEFAULT 0,
    answered_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_answers_session ON session_answers(session_id);
