package models

import (
	"time"
)

type ReminderType string

const (
	ReminderTypeOneTime    ReminderType = "one-time"
	ReminderTypePersistent ReminderType = "persistent"
	ReminderTypeRecurring  ReminderType = "recurring"
)

type NotificationChannel string

const (
	NotificationChannelEmail NotificationChannel = "email"
	NotificationChannelSMS   NotificationChannel = "sms"
	NotificationChannelCall  NotificationChannel = "call"
)

type NotificationChannels []NotificationChannel

type Reminder struct {
	ID                     int                  `json:"id"`
	UserID                 string               `json:"user_id"`
	Title                  string               `json:"title"`
	Description            *string              `json:"description"`
	ScheduledAt            time.Time            `json:"scheduled_at"`
	ReminderType           ReminderType         `json:"reminder_type"`
	NotificationChannels   NotificationChannels `json:"notification_channels"`
	ScheduledTime          time.Time            `json:"scheduled_time"`
	ScheduledDaysOfWeek    []int                `json:"scheduled_days_of_week"`
	DeliveryWindowMinutes  int                  `json:"delivery_window_minutes"`
	DeliveryMethod         string               `json:"delivery_method"`
	DeliveryAddress        string               `json:"delivery_address"`
	Status                 string               `json:"status"`
	IsActive               bool                 `json:"is_active"`
	IsPersistent           bool                 `json:"is_persistent"`
	AcknowledgedAt         *time.Time           `json:"acknowledged_at"`
	ReminderIntervalMinutes *int                `json:"reminder_interval_minutes"`
	LastRemindedAt         *time.Time           `json:"last_reminded_at"`
	CreatedAt              time.Time            `json:"created_at"`
	UpdatedAt              time.Time            `json:"updated_at"`
}

type CreateReminderRequest struct {
	Title                  string               `json:"title" validate:"required,min=1,max=500"`
	Description            *string              `json:"description"`
	ScheduledAt            time.Time            `json:"scheduled_at" validate:"required"`
	ReminderType           ReminderType         `json:"reminder_type" validate:"required,oneof=one-time persistent recurring"`
	NotificationChannels   NotificationChannels `json:"notification_channels" validate:"required,min=1"`
	ScheduledTime          string               `json:"scheduled_time" validate:"required"`
	ScheduledDaysOfWeek    []int                `json:"scheduled_days_of_week"`
	DeliveryWindowMinutes  int                  `json:"delivery_window_minutes" validate:"min=1,max=1440"`
	DeliveryMethod         string               `json:"delivery_method" validate:"required"`
	DeliveryAddress        string               `json:"delivery_address" validate:"required,email"`
	IsPersistent           bool                 `json:"is_persistent"`
	ReminderIntervalMinutes *int                `json:"reminder_interval_minutes"`
}

type UpdateReminderRequest struct {
	Title                  *string              `json:"title,omitempty"`
	Description            *string              `json:"description,omitempty"`
	ScheduledAt            *time.Time           `json:"scheduled_at,omitempty"`
	ReminderType           *ReminderType        `json:"reminder_type,omitempty"`
	NotificationChannels   NotificationChannels `json:"notification_channels,omitempty"`
	ScheduledTime          *string              `json:"scheduled_time,omitempty"`
	ScheduledDaysOfWeek    []int                `json:"scheduled_days_of_week,omitempty"`
	DeliveryWindowMinutes  *int                 `json:"delivery_window_minutes,omitempty"`
	DeliveryMethod         *string              `json:"delivery_method,omitempty"`
	DeliveryAddress        *string              `json:"delivery_address,omitempty"`
	IsActive               *bool                `json:"is_active,omitempty"`
	IsPersistent           *bool                `json:"is_persistent,omitempty"`
	ReminderIntervalMinutes *int                `json:"reminder_interval_minutes,omitempty"`
}