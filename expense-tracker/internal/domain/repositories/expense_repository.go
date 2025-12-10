package repositories

import (
	"context"
	"expense-tracker/internal/domain/entities"
	"time"
)

type ExpenseFilter struct {
	UserID    string
	StartDate *time.Time
	EndDate   *time.Time
	Category  *string
}

type ExpenseRepository interface {
	Create(ctx context.Context, expense *entities.Expense) error
	FindByID(ctx context.Context, id string) (*entities.Expense, error)
	FindByUserID(ctx context.Context, userID string, filter ExpenseFilter) ([]*entities.Expense, error)
	Update(ctx context.Context, expense *entities.Expense) error
	Delete(ctx context.Context, id string) error
	GetTotalByCategory(ctx context.Context, userID string, startDate, endDate time.Time) (map[string]float64, error)
}