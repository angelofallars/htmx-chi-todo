package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type (
	Repository interface {
		CreateUser(ctx context.Context, u *user) error
		GetUserByUsername(ctx context.Context, username username) (*user, error)
		GetUserByEmail(ctx context.Context, email email) (*user, error)
		GetUserByID(ctx context.Context, uuid uuid.UUID) (*user, error)
	}

	redisRepository struct {
		redis *redis.Client
	}
)

const (
	redisFmtUser       = "users:%v"
	redisUsernameIndex = "usersByUsername"
	redisEmailIndex    = "usersByEmail"
)

func NewRedisRepository(redis *redis.Client) Repository {
	return &redisRepository{
		redis: redis,
	}
}

func (repo redisRepository) CreateUser(ctx context.Context, u *user) (err error) {
	id := u.ID.String()

	m := map[string]any{
		"id":             id,
		"createdAt":      u.CreatedAt,
		"username":       string(u.Username),
		"email":          string(u.Email),
		"hashedPassword": string(u.HashedPassword),
	}

	pipe := repo.redis.TxPipeline()

	pipe.HSet(ctx, fmt.Sprintf(redisFmtUser, id), m)
	pipe.HSet(ctx, redisUsernameIndex,
		map[string]string{string(u.Username): id},
	)
	pipe.HSet(ctx, redisEmailIndex,
		map[string]string{string(u.Email): id},
	)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return
}

func (repo redisRepository) GetUserByEmail(ctx context.Context, email email) (*user, error) {
	cmd := repo.redis.HGet(ctx, redisEmailIndex, string(email))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	id, err := uuid.Parse(cmd.Val())
	if err != nil {
		return nil, err
	}
	return repo.GetUserByID(ctx, id)
}

func (repo redisRepository) GetUserByID(ctx context.Context, uuid uuid.UUID) (*user, error) {
	cmd := repo.redis.HGetAll(ctx, fmt.Sprintf(redisFmtUser, uuid.String()))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	u := new(user)

	err := cmd.Scan(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (repo redisRepository) GetUserByUsername(ctx context.Context, username username) (*user, error) {
	cmd := repo.redis.HGet(ctx, redisUsernameIndex, string(username))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	id, err := uuid.Parse(cmd.Val())
	if err != nil {
		return nil, err
	}
	return repo.GetUserByID(ctx, id)
}
