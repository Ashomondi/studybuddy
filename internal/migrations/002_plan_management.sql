ALTER TABLE goals ADD COLUMN current_level TEXT;
ALTER TABLE goals ADD COLUMN weekly_study_minutes INTEGER;
ALTER TABLE goals ADD COLUMN preferred_study_days TEXT NOT NULL DEFAULT '[]';

ALTER TABLE plan_topics ADD COLUMN scheduled_date TEXT;
ALTER TABLE plan_topics ADD COLUMN status TEXT NOT NULL DEFAULT 'pending'
    CHECK (status IN ('pending', 'in_progress', 'completed'));

CREATE INDEX IF NOT EXISTS idx_plan_topics_schedule ON plan_topics(scheduled_date, status);
