package services

import (
	"context"
	"errors"
	"expense-tracker/internal/application/dto"
	"expense-tracker/internal/domain/entities"
	"expense-tracker/internal/domain/repositories"

	"golang.org/x/crypto/bcrypt"
)

type JWTManager interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (string, error)
}

type AuthService struct {
	userRepo repositories.UserRepository
	jwtMgr   JWTManager
}

func NewAuthService(userRepo repositories.UserRepository, jwtMgr JWTManager) *AuthService {
	return &AuthService{userRepo: userRepo, jwtMgr: jwtMgr}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.jwtMgr.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	response := &dto.AuthResponse{
		Token: token,
	}
	response.User.ID = user.ID
	response.User.Email = user.Email
	response.User.Name = user.Name

	return response, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtMgr.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	response := &dto.AuthResponse{
		Token: token,
	}
	response.User.ID = user.ID
	response.User.Email = user.Email
	response.User.Name = user.Name

	return response, nil
}