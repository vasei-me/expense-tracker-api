package handlers

import (
	"encoding/json"
	"net/http"

	"expense-tracker/internal/application/dto"
	"expense-tracker/internal/application/services"
	"expense-tracker/internal/infrastructure/http/middleware"
	"expense-tracker/internal/pkg/validation"

	"github.com/gorilla/mux"
)

type ExpenseHandler struct {
	expenseService *services.ExpenseService
	validator      *validation.Validator
}

func NewExpenseHandler(expenseService *services.ExpenseService, validator *validation.Validator) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
		validator:      validator,
	}
}

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.expenseService.CreateExpense(r.Context(), userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *ExpenseHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var filter dto.FilterParams
	filter.Period = r.URL.Query().Get("period")
	filter.StartDate = r.URL.Query().Get("start_date")
	filter.EndDate = r.URL.Query().Get("end_date")
	filter.Category = r.URL.Query().Get("category")

	expenses, err := h.expenseService.GetExpenses(r.Context(), userID, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	expenseID := vars["id"]

	var req dto.UpdateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.expenseService.UpdateExpense(r.Context(), userID, expenseID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	expenseID := vars["id"]

	if err := h.expenseService.DeleteExpense(r.Context(), userID, expenseID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}