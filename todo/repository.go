package todo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type (
	Repository interface {
		CreateList(ctx context.Context, l *list) (*uuid.UUID, error)
		GetList(ctx context.Context, uuid *uuid.UUID) (*list, error)
		GetLists(ctx context.Context) ([]*list, error)

		GetItem(ctx context.Context, uuid *uuid.UUID) (*item, error)
		CreateItem(ctx context.Context, i *item) (*uuid.UUID, error)
		UpdateItem(ctx context.Context, uuid *uuid.UUID, i *item) error
		DeleteItem(ctx context.Context, uuid *uuid.UUID) error
	}

	redisRepository struct {
		redis *redis.Client
	}
)

const (
	redisFmtList = "lists:%v"
	redisFmtItem = "items:%v"
)

func NewRedisRepository(redis *redis.Client) Repository {
	return &redisRepository{
		redis: redis,
	}
}

func (r redisRepository) CreateList(ctx context.Context, l *list) (*uuid.UUID, error) {
	m := map[string]any{
		"id":        l.ID.String(),
		"createdAt": l.CreatedAt,
	}

	bytes, err := json.Marshal(l.Items)
	if err != nil {
		return nil, err
	}

	m["items"] = bytes

	_, err = r.redis.HSet(ctx, fmt.Sprintf(redisFmtList, l.ID.String()), m).Result()
	if err != nil {
		return nil, err
	}

	return &l.ID, nil
}

func (r redisRepository) GetList(ctx context.Context, uuid *uuid.UUID) (*list, error) {
	cmd := r.redis.HGetAll(ctx, fmt.Sprintf(redisFmtList, uuid.String()))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	l := new(list)

	err := cmd.Scan(l)
	if err != nil {
		return nil, err
	}

	items := new([]*item)

	err = json.Unmarshal([]byte(cmd.Val()["items"]), items)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (r redisRepository) GetLists(ctx context.Context) ([]*list, error) {
	// TODO
	// cmd := r.redis.HGetAll(ctx, fmt.Sprintf(redisFmtList, uuid.String()))
	// if cmd.Err() != nil {
	// 	return nil, cmd.Err()
	// }
	//
	// l := new(list)
	//
	// err := cmd.Scan(l)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// items := new([]*item)
	//
	// err = json.Unmarshal([]byte(cmd.Val()["items"]), items)
	// if err != nil {
	// 	return nil, err
	// }

	return []*list{}, nil
}

func (r redisRepository) GetItem(ctx context.Context, uuid *uuid.UUID) (*item, error) {
	cmd := r.redis.HGetAll(ctx, fmt.Sprintf(redisFmtItem, uuid.String()))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	if len(cmd.Val()) == 0 {
		return nil, errors.New("Not found")
	}

	i := new(item)

	err := cmd.Scan(i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (r redisRepository) CreateItem(ctx context.Context, i *item) (*uuid.UUID, error) {
	m := map[string]any{
		"id":          i.ID.String(),
		"createdAt":   i.CreatedAt,
		"title":       i.Title,
		"description": i.Description,
		"isDone":      bool(i.IsDone),
	}

	_, err := r.redis.HSet(ctx, fmt.Sprintf(redisFmtItem, i.ID.String()), m).Result()
	if err != nil {
		return nil, err
	}

	return &i.ID, nil
}

func (r redisRepository) UpdateItem(ctx context.Context, uuid *uuid.UUID, i *item) error {
	m := map[string]any{
		"id":          i.ID.String(),
		"createdAt":   i.CreatedAt,
		"title":       i.Title,
		"description": i.Description,
		"isDone":      bool(i.IsDone),
	}

	_, err := r.redis.HSet(ctx, fmt.Sprintf(redisFmtItem, i.ID.String()), m).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r redisRepository) DeleteItem(ctx context.Context, uuid *uuid.UUID) error {
	status, err := r.redis.Del(ctx, fmt.Sprintf(redisFmtItem, uuid.String())).Result()
	if err != nil {
		return err
	}

	hasDeletedNothing := status == 0
	if hasDeletedNothing {
		return errors.New("Deleted nothing")
	}

	return nil
}
