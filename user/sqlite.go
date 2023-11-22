package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r SQLiteRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS users(
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		createdAt INTEGER NOT NULL
	);
	`

	_, err := r.db.Exec(query)
	return err
}

func (r SQLiteRepository) CreateUser(ctx context.Context, u *User) (err error) {
	query := `INSERT INTO users( id, username, password, email, createdAt )
						  values( ?, ?, ?, ?, ? )`

	_, err = r.db.Exec(query,
		u.ID.String(),
		string(u.Username),
		string(u.HashedPassword),
		string(u.Email),
		u.CreatedAt.Unix(),
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return ErrDuplicate
			}
		}
		return err
	}

	return
}
func (r SQLiteRepository) GetUserByUsername(ctx context.Context, username Username) (*User, error) {
	query := `SELECT id, password, email, createdAt
			  FROM users
			  WHERE username = ?`

	row := r.db.QueryRow(query, string(username))

	var id string
	var password string
	var email string
	var createdAt int64
	err := row.Scan(&id, &password, &email, &createdAt)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:             uuid.MustParse(id),
		Username:       Username(username),
		HashedPassword: password,
		Email:          Email(email),
		CreatedAt:      time.Unix(createdAt, 0),
	}

	return user, nil
}

func (r SQLiteRepository) GetUserByEmail(ctx context.Context, email Email) (*User, error) {
	query := `SELECT id, password, username, createdAt
			  FROM users
			  WHERE email = ?`

	row := r.db.QueryRow(query, string(email))

	var id string
	var username string
	var password string
	var createdAt int64
	err := row.Scan(&id, &password, &email, &createdAt)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:             uuid.MustParse(id),
		Username:       Username(username),
		HashedPassword: password,
		Email:          Email(email),
		CreatedAt:      time.Unix(createdAt, 0),
	}

	return user, nil
}
func (r SQLiteRepository) GetUserByID(ctx context.Context, uuid uuid.UUID) (*User, error) {
	query := `SELECT password, username, email, createdAt
			  FROM users
			  WHERE id = ?`

	row := r.db.QueryRow(query, uuid.String())

	var username string
	var email string
	var password string
	var createdAt int64
	err := row.Scan(&username, &email, &password, &createdAt)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:             uuid,
		Username:       Username(username),
		HashedPassword: password,
		Email:          Email(email),
		CreatedAt:      time.Unix(createdAt, 0),
	}

	return user, nil
}
