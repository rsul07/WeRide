package repository

import (
	"WeRide/user-service/internal/jwt"
	"WeRide/user-service/internal/models"
	pb "WeRide/user-service/protoc/gen/go"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Repository struct {
	db       *pgxpool.Pool
	tokenTTL time.Duration
	Secret   string
}

func NewRepository(db *pgxpool.Pool, tokenTTL time.Duration, secret string) Repository {
	return Repository{
		db:       db,
		tokenTTL: tokenTTL,
		Secret:   secret,
	}
}

func (r *Repository) SaveUser(ctx context.Context, email string, password string, firstName string, lastName string, gender int64) (string, error) {
	id := uuid.New().String()
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}
	query := `
		INSERT INTO public.users (user_id, email, password_hash, first_name, last_name, gender, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.Exec(ctx, query, id, email, passHash, firstName, lastName, gender, time.Now())

	if err != nil {
		if IsUniqueViolation(err) {
			return "", errors.New("user with this email already exists")
		}
		return "", err
	}

	return id, nil
}

func (r *Repository) LoginUser(ctx context.Context, email, password string) (string, error) {

	query := `
		SELECT user_id, email, password_hash, first_name, last_name, created_at
		FROM public.users
		WHERE email = $1
		`
	row := r.db.QueryRow(ctx, query, email)
	var user models.User
	err := row.Scan(&user.UserID, &user.Email, &user.PassHash, &user.FirsName, &user.LastName, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("user not found")
		}
		return "", fmt.Errorf("error saving user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", fmt.Errorf("error checking password: %v", err)
	}

	token, err := jwt.NewToken(user, r.Secret, r.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("error creating token: %v", err)
	}
	return token, nil
}

func (r *Repository) GetUserRoutes(
	ctx context.Context,
	userID uuid.UUID,
) ([]pb.Route, error) {
	//TODO:
	result := []pb.Route{}
	return result, nil
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
