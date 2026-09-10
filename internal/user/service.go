package user

import (
	"context"
	"errors"
	"time"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/auth"
	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, req models.LoginRequest) (string, error)
	Register(ctx context.Context, req models.RegisterRequest) (string, error)
}

type service struct {
	repo Repo
	auth auth.Authenticator
	ttl  time.Duration
}

func NewService(repo Repo, a auth.Authenticator, ttl time.Duration) Service {
	return &service{repo: repo, auth: a, ttl: ttl}
}

func (s *service) Login(ctx context.Context, req models.LoginRequest) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)

	if err != nil {
		return "", err
	}

	// Hash qilingan password va user yuborgan passwordni tekshiramiz
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return "", errors.New("invalid email or password")
	}

	return s.auth.GenerateToken(user.ID, s.ttl)
}

func (s *service) Register(ctx context.Context, req models.RegisterRequest) (string, error) {
	// Passwordni hash qilamiz db ga saqlashdan oldin
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	req.Password = string(passwordHash)

	userID, err := s.repo.CreateUser(ctx, req)
	if err != nil {
		return "", err
	}

	return s.auth.GenerateToken(userID, s.ttl)
}
