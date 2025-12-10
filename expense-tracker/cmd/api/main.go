package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Config ساده
type Config struct {
	ServerPort string
	JWTSecret  string
}

func loadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // پورت پیش‌فرض
	}
	
	return &Config{
		ServerPort: port,
		JWTSecret:  "your-secret-key-change-in-production",
	}
}

// Database connection
func connectDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", "expense_tracker.db")
	if err != nil {
		return nil, err
	}
	
	// Create tables if they don't exist
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS expenses (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		amount REAL NOT NULL,
		category TEXT NOT NULL CHECK(category IN ('groceries', 'leisure', 'electronics', 'utilities', 'clothing', 'health', 'others')),
		description TEXT,
		date DATE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	
	CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
	CREATE INDEX IF NOT EXISTS idx_expenses_date ON expenses(date);
	`
	
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}
	
	log.Println("Database initialized successfully")
	return db, nil
}

// Handlers
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	response := map[string]interface{}{
		"message": "User registered successfully",
		"user": map[string]string{
			"email": req.Email,
			"name":  req.Name,
		},
		"token": "demo-jwt-token",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	response := map[string]interface{}{
		"message": "Login successful",
		"token":   "demo-jwt-token",
		"user": map[string]string{
			"email": req.Email,
			"id":    "user-123",
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createExpenseHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount      float64 `json:"amount"`
		Category    string  `json:"category"`
		Description string  `json:"description"`
		Date        string  `json:"date"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	response := map[string]interface{}{
		"message": "Expense created successfully",
		"expense": map[string]interface{}{
			"id":          "exp-123",
			"amount":      req.Amount,
			"category":    req.Category,
			"description": req.Description,
			"date":        req.Date,
			"created_at":  time.Now().Format(time.RFC3339),
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getExpensesHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering
	period := r.URL.Query().Get("period")
	category := r.URL.Query().Get("category")
	
	// Mock expenses data
	expenses := []map[string]interface{}{
		{
			"id":          "exp-1",
			"amount":      25.50,
			"category":    "groceries",
			"description": "Weekly groceries",
			"date":        time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		},
		{
			"id":          "exp-2",
			"amount":      50.00,
			"category":    "leisure",
			"description": "Movie tickets",
			"date":        time.Now().AddDate(0, 0, -3).Format("2006-01-02"),
		},
		{
			"id":          "exp-3",
			"amount":      120.00,
			"category":    "electronics",
			"description": "Phone charger",
			"date":        time.Now().AddDate(0, 0, -10).Format("2006-01-02"),
		},
	}
	
	// Apply period filter
	if period != "" {
		var filtered []map[string]interface{}
		now := time.Now()
		
		for _, expense := range expenses {
			expenseDate, err := time.Parse("2006-01-02", expense["date"].(string))
			if err != nil {
				continue
			}
			
			switch period {
			case "week":
				if expenseDate.After(now.AddDate(0, 0, -7)) {
					filtered = append(filtered, expense)
				}
			case "month":
				if expenseDate.After(now.AddDate(0, -1, 0)) {
					filtered = append(filtered, expense)
				}
			case "3months":
				if expenseDate.After(now.AddDate(0, -3, 0)) {
					filtered = append(filtered, expense)
				}
			default:
				filtered = append(filtered, expense)
			}
		}
		expenses = filtered
	}
	
	// Filter by category if provided
	if category != "" {
		filtered := []map[string]interface{}{}
		for _, expense := range expenses {
			if expense["category"] == category {
				filtered = append(filtered, expense)
			}
		}
		expenses = filtered
	}
	
	// Calculate total
	var total float64
	for _, expense := range expenses {
		total += expense["amount"].(float64)
	}
	
	response := map[string]interface{}{
		"expenses": expenses,
		"count":    len(expenses),
		"total":    total,
		"filters": map[string]string{
			"period":   period,
			"category": category,
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	expenseID := vars["id"]
	
	var req struct {
		Amount      *float64 `json:"amount"`
		Category    *string  `json:"category"`
		Description *string  `json:"description"`
		Date        *string  `json:"date"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	response := map[string]interface{}{
		"message":    "Expense updated successfully",
		"expense_id": expenseID,
		"updates":    req,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deleteExpenseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	expenseID := vars["id"]
	
	response := map[string]interface{}{
		"message":    "Expense deleted successfully",
		"expense_id": expenseID,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Auth middleware
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		
		// Simple auth check (in production, validate JWT)
		if authHeader != "Bearer demo-jwt-token" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		
		next(w, r)
	}
}

func main() {
	cfg := loadConfig()
	
	// Connect to database
	db, err := connectDB()
	if err != nil {
		log.Printf("Warning: Could not connect to database: %v", err)
		log.Println("Running in demo mode without database...")
	} else {
		defer db.Close()
		log.Println("Connected to database successfully")
	}
	
	// Initialize router
	router := mux.NewRouter()
	
	// Public routes
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Expense Tracker API v1.0"))
	})
	
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if db != nil {
			if err := db.Ping(); err != nil {
				http.Error(w, "Database connection error", http.StatusInternalServerError)
				return
			}
		}
		w.Write([]byte("OK"))
	})
	
	// Auth routes
	router.HandleFunc("/api/auth/register", registerHandler).Methods("POST")
	router.HandleFunc("/api/auth/login", loginHandler).Methods("POST")
	
	// Test endpoint
	router.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "API is working!",
			"status":  "success",
		})
	})
	
	// Protected routes
	router.HandleFunc("/api/expenses", authMiddleware(createExpenseHandler)).Methods("POST")
	router.HandleFunc("/api/expenses", authMiddleware(getExpensesHandler)).Methods("GET")
	router.HandleFunc("/api/expenses/{id}", authMiddleware(updateExpenseHandler)).Methods("PUT")
	router.HandleFunc("/api/expenses/{id}", authMiddleware(deleteExpenseHandler)).Methods("DELETE")
	
	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Println("Available endpoints:")
	log.Println("  GET  /                     - Welcome message")
	log.Println("  GET  /health               - Health check")
	log.Println("  GET  /api/test             - Test endpoint")
	log.Println("  POST /api/auth/register    - Register new user")
	log.Println("  POST /api/auth/login       - Login user")
	log.Println("  POST /api/expenses         - Create expense (protected)")
	log.Println("  GET  /api/expenses         - Get expenses (protected)")
	log.Println("  PUT  /api/expenses/{id}    - Update expense (protected)")
	log.Println("  DELETE /api/expenses/{id}  - Delete expense (protected)")
	
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}