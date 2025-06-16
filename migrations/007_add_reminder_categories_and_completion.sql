-- +goose Up
-- Add category field to reminders table
ALTER TABLE reminders 
ADD COLUMN category VARCHAR(50) NOT NULL DEFAULT 'other';

-- Add completion tracking fields
ALTER TABLE reminders 
ADD COLUMN completed_at TIMESTAMPTZ DEFAULT NULL,
ADD COLUMN completion_note TEXT DEFAULT NULL;

-- Update status constraint to include completed status
ALTER TABLE reminders DROP CONSTRAINT IF EXISTS reminders_status_check;
ALTER TABLE reminders ADD CONSTRAINT reminders_status_check 
CHECK (status IN ('pending', 'sent', 'delivered', 'acknowledged', 'completed', 'failed', 'cancelled'));

-- Create index for category-based queries
CREATE INDEX idx_reminders_category ON reminders(category);
CREATE INDEX idx_reminders_completed_at ON reminders(completed_at);

-- Create a view for statistics calculations
CREATE OR REPLACE VIEW reminder_statistics AS
WITH user_stats AS (
    SELECT 
        user_id,
        COUNT(*) as total_reminders,
        COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_reminders,
        COUNT(CASE WHEN status NOT IN ('completed', 'cancelled', 'failed') THEN 1 END) as active_reminders,
        COUNT(CASE WHEN created_at >= CURRENT_DATE - INTERVAL '7 days' THEN 1 END) as this_week_count,
        COUNT(CASE WHEN created_at >= CURRENT_DATE - INTERVAL '14 days' AND created_at < CURRENT_DATE - INTERVAL '7 days' THEN 1 END) as last_week_count,
        COUNT(CASE WHEN created_at >= CURRENT_DATE - INTERVAL '1 month' THEN 1 END) as this_month_count,
        COUNT(CASE WHEN created_at >= CURRENT_DATE - INTERVAL '2 months' AND created_at < CURRENT_DATE - INTERVAL '1 month' THEN 1 END) as last_month_count
    FROM reminders
    GROUP BY user_id
),
category_stats AS (
    SELECT 
        user_id,
        category,
        COUNT(*) as count,
        ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY user_id), 0) as percentage
    FROM reminders
    GROUP BY user_id, category
)
SELECT 
    us.user_id,
    us.total_reminders,
    us.completed_reminders,
    us.active_reminders,
    us.this_week_count,
    us.last_week_count,
    us.this_month_count,
    us.last_month_count,
    CASE 
        WHEN us.last_week_count = 0 THEN 0
        ELSE ROUND(((us.this_week_count - us.last_week_count) * 100.0 / us.last_week_count), 0)
    END as weekly_change_percent,
    CASE 
        WHEN us.last_month_count = 0 THEN 0
        ELSE ROUND(((us.this_month_count - us.last_month_count) * 100.0 / us.last_month_count), 0)
    END as monthly_change_percent,
    json_agg(
        json_build_object(
            'category', cs.category,
            'count', cs.count,
            'percentage', cs.percentage
        ) ORDER BY cs.count DESC
    ) as category_breakdown
FROM user_stats us
LEFT JOIN category_stats cs ON us.user_id = cs.user_id
GROUP BY us.user_id, us.total_reminders, us.completed_reminders, us.active_reminders, 
         us.this_week_count, us.last_week_count, us.this_month_count, us.last_month_count;

-- +goose Down
DROP VIEW IF EXISTS reminder_statistics;
DROP INDEX IF EXISTS idx_reminders_category;
DROP INDEX IF EXISTS idx_reminders_completed_at;

ALTER TABLE reminders DROP CONSTRAINT IF EXISTS reminders_status_check;
ALTER TABLE reminders ADD CONSTRAINT reminders_status_check 
CHECK (status IN ('pending', 'sent', 'delivered', 'acknowledged', 'failed', 'cancelled'));

ALTER TABLE reminders 
DROP COLUMN category,
DROP COLUMN completed_at,
DROP COLUMN completion_note;