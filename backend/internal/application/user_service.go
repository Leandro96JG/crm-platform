package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo domain.UserRepository
}

type UserService interface {
	CreateUser(ctx context.Context, user domain.User) (string, error)
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	SearchUsers(ctx context.Context, filters domain.UserFilters) (domain.PagingResult[domain.User], error)
	UpdateUser(ctx context.Context, userID string, update domain.UpdateUser) error
	DeleteUser(ctx context.Context, userID string) error
	ChangePassword(ctx context.Context, userID string, oldPassword string, newPassword string) error
	Authenticate(ctx context.Context, login string, password string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

func NewUserService(userRepo domain.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(ctx context.Context, user domain.User) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.PasswordHash = string(hash)
	return s.userRepo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *userService) SearchUsers(ctx context.Context, filters domain.UserFilters) (domain.PagingResult[domain.User], error) {
	return s.userRepo.Search(ctx, filters)
}

func (s *userService) UpdateUser(ctx context.Context, userID string, update domain.UpdateUser) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.MergeUpdate(update)
	return s.userRepo.Update(ctx, *user)
}

func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	return s.userRepo.Delete(ctx, userID)
}

func (s *userService) ChangePassword(ctx context.Context, userID string, oldPassword string, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return domain.NewUnauthorizedError("invalid password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hash)
	return s.userRepo.Update(ctx, *user)
}

func (s *userService) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

func (s *userService) Authenticate(ctx context.Context, login string, password string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, login)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = s.userRepo.FindByUsername(ctx, login)
		if err != nil {
			return nil, err
		}
	}

	if user == nil {
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}

	if !user.IsActive {
		return nil, domain.NewUnauthorizedError("user is inactive")
	}

	return user, nil
}
