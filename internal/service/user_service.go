package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"ntuc/internal/domain"
	"ntuc/internal/repository"
	"ntuc/pkg/errors"
	"ntuc/pkg/logger"
)

type UserService struct {
	repo   *repository.UserRepository
	logger *logger.Logger
}

func NewUserService(repo *repository.UserRepository, logger *logger.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

func (s *UserService) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Errorf("failed to check existing user: %v", err)
		return nil, errors.NewInternalError("failed to create user", err)
	}

	if existingUser != nil {
		return nil, errors.NewConflictError("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("failed to hash password: %v", err)
		return nil, errors.NewInternalError("failed to create user", err)
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Errorf("failed to create user: %v", err)
		return nil, errors.NewInternalError("failed to create user", err)
	}

	s.logger.Infof("user created: %s", user.ID)
	return user.ToResponse(), nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.UserResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid user id format")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Errorf("failed to get user: %v", err)
		return nil, errors.NewInternalError("failed to get user", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	return user.ToResponse(), nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]*domain.UserResponse, error) {
	users, err := s.repo.ListAll(ctx)
	if err != nil {
		s.logger.Errorf("failed to list users: %v", err)
		return nil, errors.NewInternalError("failed to list users", err)
	}

	responses := make([]*domain.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid user id format")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Errorf("failed to get user: %v", err)
		return nil, errors.NewInternalError("failed to update user", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	// Check if email is already taken by another user
	if req.Email != user.Email {
		existingUser, err := s.repo.GetByEmail(ctx, req.Email)
		if err != nil && err != sql.ErrNoRows {
			s.logger.Errorf("failed to check email: %v", err)
			return nil, errors.NewInternalError("failed to update user", err)
		}

		if existingUser != nil {
			return nil, errors.NewConflictError("email already in use")
		}
	}

	user.Name = req.Name
	user.Email = req.Email

	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Errorf("failed to update user: %v", err)
		return nil, errors.NewInternalError("failed to update user", err)
	}

	s.logger.Infof("user updated: %s", user.ID)
	return user.ToResponse(), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid user id format")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Errorf("failed to get user: %v", err)
		return errors.NewInternalError("failed to delete user", err)
	}

	if user == nil {
		return errors.NewNotFoundError("user not found")
	}

	if err := s.repo.Delete(ctx, userID); err != nil {
		s.logger.Errorf("failed to delete user: %v", err)
		return errors.NewInternalError("failed to delete user", err)
	}

	s.logger.Infof("user deleted: %s", userID)
	return nil
}
