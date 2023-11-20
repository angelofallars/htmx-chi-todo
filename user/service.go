package user

import (
	"context"
	"errors"
	"time"

	"github.com/angelofallars/htmx-chi-todo/auth"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service interface {
		Signup(ctx context.Context, u *user) error
		Login(ctx context.Context, username username, rawPassword string) (string, error)
	}

	userService struct {
		repo Repository
	}
)

func NewService(repo Repository) Service {
	return &userService{
		repo: repo,
	}
}

func (svc userService) Signup(ctx context.Context, u *user) error {
	user, err := svc.repo.GetUserByUsername(ctx, u.Username)
	if user != nil {
		return errors.New("Username already exists")
	}

	err = svc.repo.CreateUser(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

func (svc userService) Login(ctx context.Context, username username, rawPassword string) (string, error) {
	u, err := svc.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", errors.New("Incorrect email or password.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(rawPassword))
	if err != nil {
		return "", errors.New("Incorrect email or password.")
	}

	claims := &auth.JwtClaims{
		ID:       u.ID.String(),
		Username: string(u.Username),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// TODO: set up a system for getting JWT key
	signedToken, err := token.SignedString([]byte("secret"))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
