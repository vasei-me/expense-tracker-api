package repositories

import (
	"context"
	"database/sql"
	"expense-tracker/internal/domain/entities"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, email, password, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Password, user.Name,
		user.CreatedAt, user.UpdatedAt)

	return err
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, email, password, name, created_at, updated_at
		FROM users WHERE email = $1
	`

	var user entities.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.User, error) {
	query := `
		SELECT id, email, password, name, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user entities.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, email)
	return exists, err
}