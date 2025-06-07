-- +goose Up
-- Add new columns for enhanced reminder functionality
ALTER TABLE reminders 
ADD COLUMN reminder_type VARCHAR(20) NOT NULL DEFAULT 'one-time',
ADD COLUMN notification_channels JSONB NOT NULL DEFAULT '["email"]'::jsonb,
ADD COLUMN scheduled_time TIME NOT NULL DEFAULT '09:00:00',
ADD COLUMN scheduled_days_of_week INTEGER[] DEFAULT NULL,
ADD COLUMN delivery_window_minutes INTEGER NOT NULL DEFAULT 15,
ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT true;

-- Update indexes
CREATE INDEX idx_reminders_type ON reminders(reminder_type);
CREATE INDEX idx_reminders_scheduled_time ON reminders(scheduled_time);
CREATE INDEX idx_reminders_is_active ON reminders(is_active);

-- Update existing reminders to have default values
UPDATE reminders SET 
    reminder_type = 'one-time',
    notification_channels = '["email"]'::jsonb,
    scheduled_time = scheduled_at::time,
    delivery_window_minutes = 15,
    is_active = true
WHERE reminder_type IS NULL;

-- +goose Down
ALTER TABLE reminders 
DROP COLUMN reminder_type,
DROP COLUMN notification_channels,
DROP COLUMN scheduled_time,
DROP COLUMN scheduled_days_of_week,
DROP COLUMN delivery_window_minutes,
DROP COLUMN is_active;

DROP INDEX IF EXISTS idx_reminders_type;
DROP INDEX IF EXISTS idx_reminders_scheduled_time;
DROP INDEX IF EXISTS idx_reminders_is_active;