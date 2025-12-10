package entities

import (
	"expense-tracker/internal/domain/valueobjects"
	"time"
)

type Expense struct {
	ID          string                `json:"id"`
	UserID      string                `json:"user_id"`
	Amount      float64               `json:"amount"`
	Category    valueobjects.Category `json:"category"`
	Description string                `json:"description"`
	Date        time.Time             `json:"date"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}