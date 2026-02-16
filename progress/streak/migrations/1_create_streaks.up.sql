CREATE TABLE streaks (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL UNIQUE,
    current_streak   INT NOT NULL DEFAULT 0,
    longest_streak   INT NOT NULL DEFAULT 0,
    last_active_date DATE,
    daily_xp_goal    INT NOT NULL DEFAULT 20,
    streak_freezes   INT NOT NULL DEFAULT 0,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_streaks_user ON streaks(user_id);

CREATE TABLE daily_activity (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL,
    date              DATE NOT NULL,
    xp_earned         INT NOT NULL DEFAULT 0,
    lessons_completed INT NOT NULL DEFAULT 0,
    goal_met          BOOLEAN NOT NULL DEFAULT false,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, date)
);

CREATE INDEX idx_daily_activity_user ON daily_activity(user_id, date DESC);
