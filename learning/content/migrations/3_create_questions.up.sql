CREATE TABLE questions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id       UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    type            TEXT NOT NULL,              -- "listen_reply" | "multi_choice" | "single_choice"
    prompt_text     TEXT NOT NULL,              -- the question text or word to translate
    prompt_audio_key TEXT DEFAULT '',           -- object storage key for pre-recorded audio
    use_tts         BOOLEAN DEFAULT false,      -- generate audio via TTS if no recording
    correct_answer  TEXT NOT NULL,              -- correct text answer
    options         JSONB DEFAULT '[]',         -- for multi/single choice: ["option1","option2",...]
    hint            TEXT DEFAULT '',
    sort_order      INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questions_lesson ON questions(lesson_id, sort_order);
