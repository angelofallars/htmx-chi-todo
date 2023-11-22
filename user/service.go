package user

import (
	"context"
	"errors"
	"time"

	"github.com/angelofallars/htmx-chi-todo/auth"
	"github.com/angelofallars/htmx-chi-todo/service"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service interface {
		Signup(ctx context.Context, req SignupReq) error
		Login(ctx context.Context, req LoginReq) (string, error)
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

// These are the DTOs (data transfer objects)
type (
	SignupReq struct {
		Username    string
		RawPassword string
		Email       string
	}
)

type (
	LoginReq struct {
		Username    string
		RawPassword string
	}
)

func (svc userService) Signup(ctx context.Context, req SignupReq) error {
	u, err := svc.repo.GetUserByUsername(ctx, req.Username)
	if u != nil {
		return errors.New("Username already exists")
	}

	username, err := NewUsername(req.Username)
	if err != nil {
		return errors.Join(service.ErrValidation, err)
	}

	password, err := NewHashedPassword(req.RawPassword)
	if err != nil {
		return errors.Join(service.ErrValidation, err)
	}

	email, err := NewEmail(req.Email)
	if err != nil {
		return errors.Join(service.ErrValidation, err)
	}

	user := NewUser(username, password, email)
	err = svc.repo.CreateUser(ctx, user)

	if err != nil {
		return err
	}
	return nil
}

func (svc userService) Login(ctx context.Context, req LoginReq) (string, error) {
	u, err := svc.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return "", errors.New("Incorrect email or password.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(req.RawPassword))
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
