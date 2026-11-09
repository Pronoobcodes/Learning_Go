package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"

	"golang.org/x/crypto/bcrypt"
	"notes-api/internal/store/dbgen"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}

var ErrInvalidRegistration = errors.New("invalid registration")
var ErrInvalidCredentials = errors.New("invalid email or password")

type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (dbgen.User, error)
	GetUserByEmail(ctx context.Context, email string) (dbgen.User, error)
}

type AuthService struct {
	users     UserRepository
	jwtSecret []byte
}

func NewAuthService(users UserRepository, jwtSecret []byte) *AuthService {
	return &AuthService{users: users, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(
	ctx context.Context,
	email, password string,
) (dbgen.User, error) {
	email = strings.TrimSpace(email)

	if email == "" || len(password) < 8 {
		return dbgen.User{}, ErrInvalidRegistration
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return dbgen.User{}, err
	}

	return s.users.CreateUser(ctx, email, hashedPassword)
}

func (s *AuthService) Login(
	ctx context.Context,
	email, password string,
) (dbgen.User, error) {
	email = strings.TrimSpace(email)

	user, err := s.users.GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return dbgen.User{}, ErrInvalidCredentials
	}
	if err != nil {
		return dbgen.User{}, err
	}

	if CheckPassword(password, user.PasswordHash) != nil {
		return dbgen.User{}, ErrInvalidCredentials
	}

	return user, nil
}

func (s *AuthService) GenerateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID")
	}

	return int64(userID), nil
}
