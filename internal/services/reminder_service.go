package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"reminder-service/internal/db"
	"reminder-service/internal/models"
)

type ReminderService struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewReminderService(pool *pgxpool.Pool) *ReminderService {
	return &ReminderService{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (s *ReminderService) CreateReminder(ctx context.Context, userID int, req *models.CreateReminderRequest) (*models.Reminder, error) {
	scheduledTime, err := time.Parse("15:04", req.ScheduledTime)
	if err != nil {
		return nil, fmt.Errorf("invalid scheduled_time format: %w", err)
	}

	// Convert notification channels to JSON
	notificationChannelsJSON, err := json.Marshal(req.NotificationChannels)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notification channels: %w", err)
	}

	// Convert days of week to int32 slice
	var scheduledDaysOfWeek []int32
	if len(req.ScheduledDaysOfWeek) > 0 {
		scheduledDaysOfWeek = make([]int32, len(req.ScheduledDaysOfWeek))
		for i, day := range req.ScheduledDaysOfWeek {
			scheduledDaysOfWeek[i] = int32(day)
		}
	}

	params := db.CreateReminderParams{
		UserID:      int32(userID),
		Title:       req.Title,
		Description: pgtype.Text{String: "", Valid: req.Description != nil},
		ScheduledAt: pgtype.Timestamptz{
			Time:  req.ScheduledAt,
			Valid: true,
		},
		ReminderType:          string(req.ReminderType),
		NotificationChannels:  notificationChannelsJSON,
		ScheduledTime:         pgtype.Time{Microseconds: int64(scheduledTime.Hour()*3600+scheduledTime.Minute()*60) * 1000000, Valid: true},
		ScheduledDaysOfWeek:   scheduledDaysOfWeek,
		DeliveryWindowMinutes: int32(req.DeliveryWindowMinutes),
		DeliveryMethod:        req.DeliveryMethod,
		DeliveryAddress:       req.DeliveryAddress,
		IsActive:              true,
		IsPersistent:          req.IsPersistent,
		ReminderIntervalMinutes: pgtype.Int4{Valid: req.ReminderIntervalMinutes != nil},
	}

	if req.Description != nil {
		params.Description = pgtype.Text{String: *req.Description, Valid: true}
	}
	if req.ReminderIntervalMinutes != nil {
		params.ReminderIntervalMinutes = pgtype.Int4{Int32: int32(*req.ReminderIntervalMinutes), Valid: true}
	}

	reminder, err := s.queries.CreateReminder(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create reminder: %w", err)
	}

	return s.convertDBReminderToModel(&reminder), nil
}

func (s *ReminderService) GetReminder(ctx context.Context, id int32) (*models.Reminder, error) {
	reminder, err := s.queries.GetReminder(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("reminder not found")
		}
		return nil, fmt.Errorf("failed to get reminder: %w", err)
	}

	return s.convertDBReminderToModel(&reminder), nil
}

func (s *ReminderService) GetRemindersByUser(ctx context.Context, userID int) ([]models.Reminder, error) {
	dbReminders, err := s.queries.GetRemindersByUser(ctx, int32(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get reminders: %w", err)
	}

	reminders := make([]models.Reminder, len(dbReminders))
	for i, dbReminder := range dbReminders {
		reminders[i] = *s.convertDBReminderToModel(&dbReminder)
	}

	return reminders, nil
}

func (s *ReminderService) UpdateReminder(ctx context.Context, id int32, req *models.UpdateReminderRequest) (*models.Reminder, error) {
	// Get current reminder first
	current, err := s.queries.GetReminder(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("reminder not found")
		}
		return nil, fmt.Errorf("failed to get current reminder: %w", err)
	}

	params := db.UpdateReminderParams{
		ID:                    id,
		Title:                 current.Title,
		Description:           current.Description,
		ScheduledAt:           current.ScheduledAt,
		ReminderType:          current.ReminderType,
		NotificationChannels:  current.NotificationChannels,
		ScheduledTime:         current.ScheduledTime,
		ScheduledDaysOfWeek:   current.ScheduledDaysOfWeek,
		DeliveryWindowMinutes: current.DeliveryWindowMinutes,
		DeliveryMethod:        current.DeliveryMethod,
		DeliveryAddress:       current.DeliveryAddress,
		IsActive:              current.IsActive,
		IsPersistent:          current.IsPersistent,
		ReminderIntervalMinutes: current.ReminderIntervalMinutes,
	}

	// Update only provided fields
	if req.Title != nil {
		params.Title = *req.Title
	}
	if req.Description != nil {
		params.Description = pgtype.Text{String: *req.Description, Valid: true}
	}
	if req.ScheduledAt != nil {
		params.ScheduledAt = pgtype.Timestamptz{Time: *req.ScheduledAt, Valid: true}
	}
	if req.ReminderType != nil {
		params.ReminderType = string(*req.ReminderType)
	}
	if len(req.NotificationChannels) > 0 {
		notificationChannelsJSON, err := json.Marshal(req.NotificationChannels)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal notification channels: %w", err)
		}
		params.NotificationChannels = notificationChannelsJSON
	}
	if req.ScheduledTime != nil {
		scheduledTime, err := time.Parse("15:04", *req.ScheduledTime)
		if err != nil {
			return nil, fmt.Errorf("invalid scheduled_time format: %w", err)
		}
		params.ScheduledTime = pgtype.Time{Microseconds: int64(scheduledTime.Hour()*3600+scheduledTime.Minute()*60) * 1000000, Valid: true}
	}
	if len(req.ScheduledDaysOfWeek) > 0 {
		scheduledDaysOfWeek := make([]int32, len(req.ScheduledDaysOfWeek))
		for i, day := range req.ScheduledDaysOfWeek {
			scheduledDaysOfWeek[i] = int32(day)
		}
		params.ScheduledDaysOfWeek = scheduledDaysOfWeek
	}
	if req.DeliveryWindowMinutes != nil {
		params.DeliveryWindowMinutes = int32(*req.DeliveryWindowMinutes)
	}
	if req.DeliveryMethod != nil {
		params.DeliveryMethod = *req.DeliveryMethod
	}
	if req.DeliveryAddress != nil {
		params.DeliveryAddress = *req.DeliveryAddress
	}
	if req.IsActive != nil {
		params.IsActive = *req.IsActive
	}
	if req.IsPersistent != nil {
		params.IsPersistent = *req.IsPersistent
	}
	if req.ReminderIntervalMinutes != nil {
		params.ReminderIntervalMinutes = pgtype.Int4{Int32: int32(*req.ReminderIntervalMinutes), Valid: true}
	}

	reminder, err := s.queries.UpdateReminder(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update reminder: %w", err)
	}

	return s.convertDBReminderToModel(&reminder), nil
}

func (s *ReminderService) DeleteReminder(ctx context.Context, id int32) error {
	err := s.queries.DeleteReminder(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete reminder: %w", err)
	}
	return nil
}

func (s *ReminderService) GetRemindersDueForDelivery(ctx context.Context) ([]models.Reminder, error) {
	dbReminders, err := s.queries.GetRemindersDueForDelivery(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get reminders due for delivery: %w", err)
	}

	reminders := make([]models.Reminder, len(dbReminders))
	for i, dbReminder := range dbReminders {
		reminders[i] = *s.convertDBReminderToModel(&dbReminder)
	}

	return reminders, nil
}

func (s *ReminderService) convertDBReminderToModel(dbReminder *db.Reminder) *models.Reminder {
	reminder := &models.Reminder{
		ID:                    int(dbReminder.ID),
		UserID:                int(dbReminder.UserID),
		Title:                 dbReminder.Title,
		ReminderType:          models.ReminderType(dbReminder.ReminderType),
		DeliveryWindowMinutes: int(dbReminder.DeliveryWindowMinutes),
		DeliveryMethod:        dbReminder.DeliveryMethod,
		DeliveryAddress:       dbReminder.DeliveryAddress,
		Status:                dbReminder.Status,
		IsActive:              dbReminder.IsActive,
		IsPersistent:          dbReminder.IsPersistent,
	}

	// Handle nullable description
	if dbReminder.Description.Valid {
		reminder.Description = &dbReminder.Description.String
	}

	// Handle scheduled_at
	if dbReminder.ScheduledAt.Valid {
		reminder.ScheduledAt = dbReminder.ScheduledAt.Time
	}

	// Handle scheduled_time
	if dbReminder.ScheduledTime.Valid {
		// Convert microseconds to time
		totalSeconds := dbReminder.ScheduledTime.Microseconds / 1000000
		hours := totalSeconds / 3600
		minutes := (totalSeconds % 3600) / 60
		seconds := totalSeconds % 60
		reminder.ScheduledTime = time.Date(0, 1, 1, int(hours), int(minutes), int(seconds), 0, time.UTC)
	}

	// Handle created_at and updated_at
	if dbReminder.CreatedAt.Valid {
		reminder.CreatedAt = dbReminder.CreatedAt.Time
	}
	if dbReminder.UpdatedAt.Valid {
		reminder.UpdatedAt = dbReminder.UpdatedAt.Time
	}

	// Handle notification channels JSON
	if len(dbReminder.NotificationChannels) > 0 {
		var channels models.NotificationChannels
		if err := json.Unmarshal(dbReminder.NotificationChannels, &channels); err == nil {
			reminder.NotificationChannels = channels
		}
	}

	// Handle scheduled days of week
	if len(dbReminder.ScheduledDaysOfWeek) > 0 {
		reminder.ScheduledDaysOfWeek = make([]int, len(dbReminder.ScheduledDaysOfWeek))
		for i, day := range dbReminder.ScheduledDaysOfWeek {
			reminder.ScheduledDaysOfWeek[i] = int(day)
		}
	}

	// Handle acknowledged_at
	if dbReminder.AcknowledgedAt.Valid {
		reminder.AcknowledgedAt = &dbReminder.AcknowledgedAt.Time
	}

	// Handle reminder_interval_minutes
	if dbReminder.ReminderIntervalMinutes.Valid {
		intervalMinutes := int(dbReminder.ReminderIntervalMinutes.Int32)
		reminder.ReminderIntervalMinutes = &intervalMinutes
	}

	// Handle last_reminded_at
	if dbReminder.LastRemindedAt.Valid {
		reminder.LastRemindedAt = &dbReminder.LastRemindedAt.Time
	}

	return reminder
}

func (s *ReminderService) AcknowledgeReminder(ctx context.Context, id int32) error {
	err := s.queries.AcknowledgeReminder(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to acknowledge reminder: %w", err)
	}
	return nil
}

func (s *ReminderService) UpdateLastRemindedAt(ctx context.Context, id int32) error {
	err := s.queries.UpdateLastRemindedAt(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to update last reminded at: %w", err)
	}
	return nil
}

func (s *ReminderService) GetUnacknowledgedPersistentReminders(ctx context.Context) ([]models.Reminder, error) {
	dbReminders, err := s.queries.GetUnacknowledgedPersistentReminders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get unacknowledged persistent reminders: %w", err)
	}

	reminders := make([]models.Reminder, len(dbReminders))
	for i, dbReminder := range dbReminders {
		reminders[i] = *s.convertDBReminderToModel(&dbReminder)
	}

	return reminders, nil
}

func (s *ReminderService) GetPendingRecurringReminders(ctx context.Context) ([]models.Reminder, error) {
	dbReminders, err := s.queries.GetPendingRecurringReminders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending recurring reminders: %w", err)
	}

	reminders := make([]models.Reminder, len(dbReminders))
	for i, dbReminder := range dbReminders {
		reminders[i] = *s.convertDBReminderToModel(&dbReminder)
	}

	return reminders, nil
}