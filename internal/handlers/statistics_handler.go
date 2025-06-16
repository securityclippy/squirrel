package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"reminder-service/internal/middleware"
	"reminder-service/internal/services"
)

type StatisticsHandler struct {
	statisticsService *services.StatisticsService
}

func NewStatisticsHandler(statisticsService *services.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHandler) RegisterRoutes(router *mux.Router, authMiddleware *middleware.AuthMiddleware) {
	// Statistics routes
	statsRouter := router.PathPrefix("/api/statistics").Subrouter()
	// TODO: Re-enable auth middleware when ready
	// statsRouter.Use(authMiddleware.AuthenticateAny)
	
	statsRouter.HandleFunc("", h.GetUserStatistics).Methods("GET")
	statsRouter.HandleFunc("/complete/{id}", h.CompleteReminder).Methods("POST")
	statsRouter.HandleFunc("/category/{id}", h.UpdateReminderCategory).Methods("PUT")
}

// GetUserStatistics handles GET /api/statistics
func (h *StatisticsHandler) GetUserStatistics(w http.ResponseWriter, r *http.Request) {
	// TODO: Get user ID from auth context when auth is re-enabled
	// For now, use a default user ID for testing
	userID := int32(1)
	
	// Get user ID from query parameter for testing
	if userIDParam := r.URL.Query().Get("user_id"); userIDParam != "" {
		if id, err := strconv.ParseInt(userIDParam, 10, 32); err == nil {
			userID = int32(id)
		}
	}

	stats, err := h.statisticsService.GetUserStatistics(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch statistics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Failed to encode statistics", http.StatusInternalServerError)
		return
	}
}

// CompleteReminder handles POST /api/statistics/complete/{id}
func (h *StatisticsHandler) CompleteReminder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reminderIDStr := vars["id"]
	
	reminderID, err := strconv.ParseInt(reminderIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid reminder ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Note string `json:"note"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.statisticsService.CompleteReminder(r.Context(), int32(reminderID), request.Note)
	if err != nil {
		http.Error(w, "Failed to complete reminder: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Reminder marked as completed",
	})
}

// UpdateReminderCategory handles PUT /api/statistics/category/{id}
func (h *StatisticsHandler) UpdateReminderCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reminderIDStr := vars["id"]
	
	reminderID, err := strconv.ParseInt(reminderIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid reminder ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Category string `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate category
	validCategories := map[string]bool{
		"work":     true,
		"personal": true,
		"health":   true,
		"other":    true,
	}

	if !validCategories[request.Category] {
		http.Error(w, "Invalid category. Must be one of: work, personal, health, other", http.StatusBadRequest)
		return
	}

	err = h.statisticsService.UpdateReminderCategory(r.Context(), int32(reminderID), request.Category)
	if err != nil {
		http.Error(w, "Failed to update reminder category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Reminder category updated",
	})
}