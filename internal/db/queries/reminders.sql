-- name: CreateReminder :one
INSERT INTO reminders (
    user_id, title, description, scheduled_at, reminder_type, 
    notification_channels, scheduled_time, scheduled_days_of_week, 
    delivery_window_minutes, delivery_method, delivery_address, is_active,
    is_persistent, reminder_interval_minutes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: GetReminder :one
SELECT * FROM reminders WHERE id = $1;

-- name: GetRemindersByUser :many
SELECT * FROM reminders 
WHERE user_id = $1 AND is_active = true 
ORDER BY scheduled_at ASC;

-- name: GetActiveReminders :many
SELECT * FROM reminders 
WHERE is_active = true AND status = 'pending'
ORDER BY scheduled_at ASC;

-- name: GetRemindersDueForDelivery :many
SELECT * FROM reminders 
WHERE is_active = true 
  AND status = 'pending'
  AND scheduled_at <= NOW() + INTERVAL '1 minute' * delivery_window_minutes
  AND scheduled_at >= NOW() - INTERVAL '1 minute' * delivery_window_minutes;

-- name: UpdateReminder :one
UPDATE reminders SET
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    scheduled_at = COALESCE($4, scheduled_at),
    reminder_type = COALESCE($5, reminder_type),
    notification_channels = COALESCE($6, notification_channels),
    scheduled_time = COALESCE($7, scheduled_time),
    scheduled_days_of_week = COALESCE($8, scheduled_days_of_week),
    delivery_window_minutes = COALESCE($9, delivery_window_minutes),
    delivery_method = COALESCE($10, delivery_method),
    delivery_address = COALESCE($11, delivery_address),
    is_active = COALESCE($12, is_active),
    is_persistent = COALESCE($13, is_persistent),
    reminder_interval_minutes = COALESCE($14, reminder_interval_minutes),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateReminderStatus :exec
UPDATE reminders SET 
    status = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteReminder :exec
UPDATE reminders SET 
    is_active = false,
    updated_at = NOW()
WHERE id = $1;

-- name: GetRecurringReminders :many
SELECT * FROM reminders 
WHERE reminder_type IN ('persistent', 'recurring')
  AND is_active = true
  AND (
    EXTRACT(DOW FROM NOW()) = ANY(scheduled_days_of_week)
    OR scheduled_days_of_week IS NULL
  );

-- name: AcknowledgeReminder :exec
UPDATE reminders SET 
    acknowledged_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateLastRemindedAt :exec
UPDATE reminders SET 
    last_reminded_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: GetUnacknowledgedPersistentReminders :many
SELECT * FROM reminders 
WHERE is_persistent = true 
  AND is_active = true
  AND acknowledged_at IS NULL
  AND (
    last_reminded_at IS NULL 
    OR last_reminded_at < NOW() - INTERVAL '1 minute' * reminder_interval_minutes
  );

-- name: GetPendingRecurringReminders :many
SELECT * FROM reminders 
WHERE reminder_type = 'recurring'
  AND is_active = true
  AND status = 'pending'
  AND scheduled_time <= NOW()::time
  AND (
    scheduled_days_of_week IS NULL
    OR EXTRACT(DOW FROM NOW()) = ANY(scheduled_days_of_week)
  );