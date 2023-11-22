package user

import (
	"errors"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uuid.UUID      `redis:"id"`
	CreatedAt      time.Time      `redis:"createdAt"`
	Username       Username       `redis:"username"`
	Email          Email          `redis:"email"`
	HashedPassword HashedPassword `redis:"hashedPassword"`
}

func NewUser(username Username, hashedPassword HashedPassword, email Email) *User {
	return &User{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		Username:       username,
		HashedPassword: hashedPassword,
		Email:          email,
	}
}

type Username = string

func NewUsername(un string) (Username, error) {
	errs := make([]error, 0, 3)

	if len(un) < 4 {
		errs = append(errs, errors.New("Username must be at least 4 characters"))
	}

	if len(un) > 16 {
		errs = append(errs, errors.New("Username must be less than 16 characters"))
	}

	for _, r := range un {
		isNotLetter := !unicode.IsLetter(r)
		isNotNumber := !unicode.IsNumber(r)
		isNotAlphanumeric := isNotLetter && isNotNumber

		if isNotAlphanumeric {
			errs = append(errs, errors.New("Username must only contain letters, numbers and underscores"))
			break
		}
	}

	if len(errs) != 0 {
		return Username(""), errors.Join(errs...)
	}

	return Username(un), nil
}

type Email = string

func NewEmail(e string) (Email, error) {
	if err := validator.New().Var(e, "email"); err != nil {
		return Email(""), err
	}

	return Email(e), nil
}

type HashedPassword = string

func NewHashedPassword(rawPassword string) (HashedPassword, error) {
	errs := make([]error, 0, 1)

	if len(rawPassword) < 8 {
		errs = append(errs, errors.New("Username must be at least 8 characters"))
	}

	hp, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return HashedPassword(""), errors.Join(errs...)
	}

	return HashedPassword(string(hp)), nil
}
