-- Daily review goal (used by stats progress)
ALTER TABLE users ADD COLUMN daily_goal INTEGER NOT NULL DEFAULT 20 CHECK (daily_goal > 0);
