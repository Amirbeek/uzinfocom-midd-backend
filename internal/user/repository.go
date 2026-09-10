package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/database"
	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailTaken   = errors.New("email already registered")
)

type Repo interface {
	CreateUser(ctx context.Context, order models.RegisterRequest) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type repo struct{ db *sql.DB }

func NewRepo(db database.Service) Repo {
	return &repo{db: db.DB()}
}

func (r *repo) CreateUser(ctx context.Context, req models.RegisterRequest) (int64, error) {

	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`, req.Name, req.Email, req.Password).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	err := r.db.QueryRow(`
		SELECT id, name, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
