package user

import (
	"errors"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID             uuid.UUID      `redis:"id"`
	CreatedAt      time.Time      `redis:"createdAt"`
	Username       username       `redis:"username"`
	Email          email          `redis:"email"`
	HashedPassword hashedPassword `redis:"hashedPassword"`
}

func newUser(username username, hashedPassword hashedPassword, email email) *user {
	return &user{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		Username:       username,
		HashedPassword: hashedPassword,
		Email:          email,
	}
}

type username = string

func newUsername(un string) (username, error) {
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
		return username(""), errors.Join(errs...)
	}

	return username(un), nil
}

type email = string

func newEmail(e string) (email, error) {
	if err := validator.New().Var(e, "email"); err != nil {
		return email(""), err
	}

	return email(e), nil
}

type hashedPassword = string

func newHashedPassword(rawPassword string) (hashedPassword, error) {
	errs := make([]error, 0, 1)

	if len(rawPassword) < 8 {
		errs = append(errs, errors.New("Username must be at least 8 characters"))
	}

	hp, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return hashedPassword(""), errors.Join(errs...)
	}

	return hashedPassword(string(hp)), nil
}
