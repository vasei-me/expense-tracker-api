package services

import (
	"context"
	"errors"
	"expense-tracker/internal/application/dto"
	"expense-tracker/internal/domain/entities"
	"expense-tracker/internal/domain/repositories"
	"expense-tracker/internal/domain/valueobjects"
	"time"
)

type ExpenseService struct {
	expenseRepo repositories.ExpenseRepository
}

func NewExpenseService(expenseRepo repositories.ExpenseRepository) *ExpenseService {
	return &ExpenseService{expenseRepo: expenseRepo}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, userID string, req dto.CreateExpenseRequest) (*dto.ExpenseResponse, error) {
	category := valueobjects.Category(req.Category)
	if !category.IsValid() {
		return nil, errors.New("invalid category")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	expense := &entities.Expense{
		UserID:      userID,
		Amount:      req.Amount,
		Category:    category,
		Description: req.Description,
		Date:        date,
	}

	if err := s.expenseRepo.Create(ctx, expense); err != nil {
		return nil, err
	}

	return s.toResponse(expense), nil
}

func (s *ExpenseService) GetExpenses(ctx context.Context, userID string, filter dto.FilterParams) ([]*dto.ExpenseResponse, error) {
	expenseFilter := repositories.ExpenseFilter{
		UserID: userID,
	}

	if filter.Category != "" {
		expenseFilter.Category = &filter.Category
	}

	// Apply period filters
	now := time.Now()
	switch filter.Period {
	case "week":
		startDate := now.AddDate(0, 0, -7)
		expenseFilter.StartDate = &startDate
		expenseFilter.EndDate = &now
	case "month":
		startDate := now.AddDate(0, -1, 0)
		expenseFilter.StartDate = &startDate
		expenseFilter.EndDate = &now
	case "3months":
		startDate := now.AddDate(0, -3, 0)
		expenseFilter.StartDate = &startDate
		expenseFilter.EndDate = &now
	case "custom":
		if filter.StartDate != "" && filter.EndDate != "" {
			startDate, err1 := time.Parse("2006-01-02", filter.StartDate)
			endDate, err2 := time.Parse("2006-01-02", filter.EndDate)
			if err1 == nil && err2 == nil {
				expenseFilter.StartDate = &startDate
				expenseFilter.EndDate = &endDate
			}
		}
	}

	expenses, err := s.expenseRepo.FindByUserID(ctx, userID, expenseFilter)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		responses[i] = s.toResponse(expense)
	}

	return responses, nil
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, userID, expenseID string, req dto.UpdateExpenseRequest) (*dto.ExpenseResponse, error) {
	expense, err := s.expenseRepo.FindByID(ctx, expenseID)
	if err != nil {
		return nil, errors.New("expense not found")
	}

	if expense.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if req.Amount != nil {
		expense.Amount = *req.Amount
	}
	if req.Category != nil {
		category := valueobjects.Category(*req.Category)
		if !category.IsValid() {
			return nil, errors.New("invalid category")
		}
		expense.Category = category
	}
	if req.Description != nil {
		expense.Description = *req.Description
	}
	if req.Date != nil {
		date, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, errors.New("invalid date format")
		}
		expense.Date = date
	}

	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return nil, err
	}

	return s.toResponse(expense), nil
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, userID, expenseID string) error {
	expense, err := s.expenseRepo.FindByID(ctx, expenseID)
	if err != nil {
		return errors.New("expense not found")
	}

	if expense.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.expenseRepo.Delete(ctx, expenseID)
}

func (s *ExpenseService) toResponse(expense *entities.Expense) *dto.ExpenseResponse {
	return &dto.ExpenseResponse{
		ID:          expense.ID,
		Amount:      expense.Amount,
		Category:    string(expense.Category),
		Description: expense.Description,
		Date:        expense.Date,
		CreatedAt:   expense.CreatedAt,
		UpdatedAt:   expense.UpdatedAt,
	}
}