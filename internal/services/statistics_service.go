package services

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"reminder-service/internal/db"
)

type StatisticsService struct {
	queries *db.Queries
}

type OverviewStats struct {
	TotalReminders     int64     `json:"total_reminders"`
	CompletedReminders int64     `json:"completed_reminders"`
	ActiveReminders    int64     `json:"active_reminders"`
	MemberSince        time.Time `json:"member_since"`
}

type WeeklyStats struct {
	ThisWeek      int64 `json:"this_week"`
	LastWeek      int64 `json:"last_week"`
	ChangePercent int64 `json:"change_percent"`
}

type MonthlyStats struct {
	ThisMonth     int64 `json:"this_month"`
	LastMonth     int64 `json:"last_month"`
	ChangePercent int64 `json:"change_percent"`
}

type CategoryStat struct {
	Name       string `json:"name"`
	Count      int64  `json:"count"`
	Percentage int64  `json:"percentage"`
}

type UserStatistics struct {
	Overview   OverviewStats   `json:"overview"`
	Weekly     WeeklyStats     `json:"weekly"`
	Monthly    MonthlyStats    `json:"monthly"`
	Categories []CategoryStat  `json:"categories"`
}

func NewStatisticsService(queries *db.Queries) *StatisticsService {
	return &StatisticsService{
		queries: queries,
	}
}

func (s *StatisticsService) GetUserStatistics(ctx context.Context, userID int32) (*UserStatistics, error) {
	// Get overview stats
	overviewResult, err := s.queries.GetUserOverviewStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get weekly stats
	weeklyResult, err := s.queries.GetUserWeeklyStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get monthly stats
	monthlyResult, err := s.queries.GetUserMonthlyStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get category stats
	categoryResults, err := s.queries.GetUserCategoryStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert member_since from interface{} to time.Time
	var memberSince time.Time
	if overviewResult.MemberSince != nil {
		if ts, ok := overviewResult.MemberSince.(time.Time); ok {
			memberSince = ts
		}
	}

	// Calculate percentage changes
	weeklyChangePercent := calculatePercentageChange(weeklyResult.ThisWeek, weeklyResult.LastWeek)
	monthlyChangePercent := calculatePercentageChange(monthlyResult.ThisMonth, monthlyResult.LastMonth)

	// Convert category results
	categories := make([]CategoryStat, len(categoryResults))
	for i, cat := range categoryResults {
		var percentage int64
		if cat.Percentage.Valid {
			if val, err := cat.Percentage.Int64Value(); err == nil {
				percentage = val.Int64
			}
		}
		
		categories[i] = CategoryStat{
			Name:       cat.Category,
			Count:      cat.Count,
			Percentage: percentage,
		}
	}

	return &UserStatistics{
		Overview: OverviewStats{
			TotalReminders:     overviewResult.TotalReminders,
			CompletedReminders: overviewResult.CompletedReminders,
			ActiveReminders:    overviewResult.ActiveReminders,
			MemberSince:        memberSince,
		},
		Weekly: WeeklyStats{
			ThisWeek:      weeklyResult.ThisWeek,
			LastWeek:      weeklyResult.LastWeek,
			ChangePercent: weeklyChangePercent,
		},
		Monthly: MonthlyStats{
			ThisMonth:     monthlyResult.ThisMonth,
			LastMonth:     monthlyResult.LastMonth,
			ChangePercent: monthlyChangePercent,
		},
		Categories: categories,
	}, nil
}

func (s *StatisticsService) CompleteReminder(ctx context.Context, reminderID int32, note string) error {
	var pgNote pgtype.Text
	if note != "" {
		pgNote = pgtype.Text{String: note, Valid: true}
	}
	
	return s.queries.CompleteReminder(ctx, db.CompleteReminderParams{
		ID:             reminderID,
		CompletionNote: pgNote,
	})
}

func (s *StatisticsService) UpdateReminderCategory(ctx context.Context, reminderID int32, category string) error {
	return s.queries.UpdateReminderCategory(ctx, db.UpdateReminderCategoryParams{
		ID:       reminderID,
		Category: category,
	})
}

// calculatePercentageChange calculates the percentage change between current and previous values
func calculatePercentageChange(current, previous int64) int64 {
	if previous == 0 {
		if current > 0 {
			return 100 // If previous was 0 and current > 0, it's a 100% increase
		}
		return 0
	}
	
	change := ((current - previous) * 100) / previous
	return change
}