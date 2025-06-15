package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"reminder-service/internal/middleware"
	"reminder-service/internal/models"
	"reminder-service/internal/services"
)

type ReminderHandler struct {
	reminderService *services.ReminderService
}

func NewReminderHandler(reminderService *services.ReminderService) *ReminderHandler {
	return &ReminderHandler{
		reminderService: reminderService,
	}
}

func (h *ReminderHandler) CreateReminder(w http.ResponseWriter, r *http.Request) {
	// Temporarily use default user ID 1 for server-side calls
	userID := 1

	var req models.CreateReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	reminder, err := h.reminderService.CreateReminder(r.Context(), userID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reminder)
}

func (h *ReminderHandler) GetReminder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid reminder ID", http.StatusBadRequest)
		return
	}

	reminder, err := h.reminderService.GetReminder(r.Context(), int32(id))
	if err != nil {
		if err.Error() == "reminder not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminder)
}

func (h *ReminderHandler) GetReminders(w http.ResponseWriter, r *http.Request) {
	// Temporarily use default user ID 1 for server-side calls
	userID := 1

	reminders, err := h.reminderService.GetRemindersByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"reminders": reminders,
		"count":     len(reminders),
	})
}

func (h *ReminderHandler) UpdateReminder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid reminder ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	reminder, err := h.reminderService.UpdateReminder(r.Context(), int32(id), &req)
	if err != nil {
		if err.Error() == "reminder not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminder)
}

func (h *ReminderHandler) DeleteReminder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid reminder ID", http.StatusBadRequest)
		return
	}

	err = h.reminderService.DeleteReminder(r.Context(), int32(id))
	if err != nil {
		if err.Error() == "reminder not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ReminderHandler) RegisterRoutes(router *mux.Router, authMiddleware *middleware.AuthMiddleware) {
	// Reminder routes - temporarily disable auth for server-side calls
	reminderRouter := router.PathPrefix("/api/reminders").Subrouter()
	// reminderRouter.Use(authMiddleware.AuthenticateAny)

	reminderRouter.HandleFunc("", h.CreateReminder).Methods("POST")
	reminderRouter.HandleFunc("", h.GetReminders).Methods("GET")
	reminderRouter.HandleFunc("/{id:[0-9]+}", h.GetReminder).Methods("GET")
	reminderRouter.HandleFunc("/{id:[0-9]+}", h.UpdateReminder).Methods("PUT")
	reminderRouter.HandleFunc("/{id:[0-9]+}", h.DeleteReminder).Methods("DELETE")
}