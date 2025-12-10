package repositories

import (
	"context"
	"database/sql"
	"expense-tracker/internal/domain/entities"
	"expense-tracker/internal/domain/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ExpenseRepositoryImpl struct {
	db *sqlx.DB
}

func NewExpenseRepository(db *sqlx.DB) *ExpenseRepositoryImpl {
	return &ExpenseRepositoryImpl{db: db}
}

func (r *ExpenseRepositoryImpl) Create(ctx context.Context, expense *entities.Expense) error {
	expense.ID = uuid.New().String()
	expense.CreatedAt = time.Now()
	expense.UpdatedAt = time.Now()

	query := `
		INSERT INTO expenses (id, user_id, amount, category, description, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		expense.ID, expense.UserID, expense.Amount, expense.Category,
		expense.Description, expense.Date, expense.CreatedAt, expense.UpdatedAt)

	return err
}

func (r *ExpenseRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.Expense, error) {
	query := `
		SELECT id, user_id, amount, category, description, date, created_at, updated_at
		FROM expenses WHERE id = $1
	`

	var expense entities.Expense
	err := r.db.GetContext(ctx, &expense, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &expense, err
}

func (r *ExpenseRepositoryImpl) FindByUserID(ctx context.Context, userID string, filter repositories.ExpenseFilter) ([]*entities.Expense, error) {
	query := `SELECT id, user_id, amount, category, description, date, created_at, updated_at FROM expenses WHERE user_id = $1`
	args := []interface{}{userID}
	argIndex := 2

	if filter.StartDate != nil {
		query += ` AND date >= $` + fmt.Sprintf("%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		query += ` AND date <= $` + fmt.Sprintf("%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	if filter.Category != nil {
		query += ` AND category = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, *filter.Category)
		argIndex++
	}

	query += ` ORDER BY date DESC`

	var expenses []*entities.Expense
	err := r.db.SelectContext(ctx, &expenses, query, args...)
	if err == sql.ErrNoRows {
		return []*entities.Expense{}, nil
	}
	return expenses, err
}

func (r *ExpenseRepositoryImpl) Update(ctx context.Context, expense *entities.Expense) error {
	expense.UpdatedAt = time.Now()

	query := `
		UPDATE expenses 
		SET amount = $1, category = $2, description = $3, date = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.ExecContext(ctx, query,
		expense.Amount, expense.Category, expense.Description,
		expense.Date, expense.UpdatedAt, expense.ID)

	return err
}

func (r *ExpenseRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM expenses WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ExpenseRepositoryImpl) GetTotalByCategory(ctx context.Context, userID string, startDate, endDate time.Time) (map[string]float64, error) {
	query := `
		SELECT category, SUM(amount) as total
		FROM expenses 
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		GROUP BY category
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]float64)
	for rows.Next() {
		var category string
		var total float64
		if err := rows.Scan(&category, &total); err != nil {
			return nil, err
		}
		result[category] = total
	}

	return result, nil
}