package service

import (
	"context"
	"errors"

	"github.com/auhmaugmaufm/event-driven-order/internal/auth"
	"github.com/auhmaugmaufm/event-driven-order/internal/domain"
	"github.com/auhmaugmaufm/event-driven-order/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo       domain.UserRepository
	jwtManager *auth.JWTManager
}

func NewUserService(repo domain.UserRepository, jwtManager *auth.JWTManager) *UserService {
	return &UserService{repo: repo, jwtManager: jwtManager}
}
func (s *UserService) Create(ctx context.Context, req *dto.UserRequest) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(bytes),
	}
	return s.repo.Create(user)
}

func (s *UserService) Login(email string, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}
	return s.jwtManager.GenerateToken(user.ID, user.Email)
}
