package interfaces

import "net/http"

type ExpenseHandler interface {
	CreateExpense(w http.ResponseWriter, r *http.Request)
	GetExpenses(w http.ResponseWriter, r *http.Request)
	UpdateExpense(w http.ResponseWriter, r *http.Request)
	DeleteExpense(w http.ResponseWriter, r *http.Request)
}