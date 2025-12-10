package dto

import (
	"time"
)

type CreateExpenseRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Category    string  `json:"category" validate:"required"`
	Description string  `json:"description" validate:"max=500"`
	Date        string  `json:"date" validate:"required,datetime=2006-01-02"`
}

type UpdateExpenseRequest struct {
	Amount      *float64 `json:"amount" validate:"omitempty,gt=0"`
	Category    *string  `json:"category" validate:"omitempty"`
	Description *string  `json:"description" validate:"omitempty,max=500"`
	Date        *string  `json:"date" validate:"omitempty,datetime=2006-01-02"`
}

type ExpenseResponse struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FilterParams struct {
	Period     string `query:"period"` // week, month, 3months, custom
	StartDate  string `query:"start_date"`
	EndDate    string `query:"end_date"`
	Category   string `query:"category"`
}