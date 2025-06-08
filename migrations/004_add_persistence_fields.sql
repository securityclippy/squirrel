-- +goose Up
-- Add fields for persistence and acknowledgment functionality
ALTER TABLE reminders 
ADD COLUMN is_persistent BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN acknowledged_at TIMESTAMPTZ DEFAULT NULL,
ADD COLUMN reminder_interval_minutes INTEGER DEFAULT NULL,
ADD COLUMN last_reminded_at TIMESTAMPTZ DEFAULT NULL;

-- Update reminder_type to support recurring reminders
ALTER TABLE reminders DROP CONSTRAINT IF EXISTS reminders_reminder_type_check;
ALTER TABLE reminders ADD CONSTRAINT reminders_reminder_type_check 
CHECK (reminder_type IN ('one-time', 'persistent', 'recurring'));

-- Add indexes for new fields
CREATE INDEX idx_reminders_is_persistent ON reminders(is_persistent);
CREATE INDEX idx_reminders_acknowledged_at ON reminders(acknowledged_at);
CREATE INDEX idx_reminders_last_reminded_at ON reminders(last_reminded_at);

-- +goose Down
ALTER TABLE reminders 
DROP COLUMN is_persistent,
DROP COLUMN acknowledged_at,
DROP COLUMN reminder_interval_minutes,
DROP COLUMN last_reminded_at;

ALTER TABLE reminders DROP CONSTRAINT IF EXISTS reminders_reminder_type_check;
ALTER TABLE reminders ADD CONSTRAINT reminders_reminder_type_check 
CHECK (reminder_type IN ('one-time', 'persistent'));

DROP INDEX IF EXISTS idx_reminders_is_persistent;
DROP INDEX IF EXISTS idx_reminders_acknowledged_at;
DROP INDEX IF EXISTS idx_reminders_last_reminded_at;