package todo

import (
	"context"

	"github.com/google/uuid"
)

type (
	Service interface {
		CreateList(ctx context.Context, l *list) (*uuid.UUID, error)
		GetList(ctx context.Context, id *uuid.UUID) (*list, error)
		GetLists(ctx context.Context) ([]*list, error)
		GetItem(ctx context.Context, id *uuid.UUID) (*item, error)
		CreateItem(ctx context.Context, i *item) (*uuid.UUID, error)
		UpdateItem(ctx context.Context, id *uuid.UUID, title string, description string) (*item, error)
		ToggleItemComplete(ctx context.Context, id *uuid.UUID) (isDone, error)
		DeleteItem(ctx context.Context, id *uuid.UUID) error
	}

	service struct {
		repo Repository
	}
)

func NewService(repository Repository) Service {
	return service{
		repo: repository,
	}
}

func (sv service) CreateList(ctx context.Context, l *list) (*uuid.UUID, error) {
	return sv.repo.CreateList(ctx, l)
}

func (sv service) GetList(ctx context.Context, id *uuid.UUID) (*list, error) {
	return sv.repo.GetList(ctx, id)
}

func (sv service) GetLists(ctx context.Context) ([]*list, error) {
	return sv.repo.GetLists(ctx)
}

func (sv service) GetItem(ctx context.Context, id *uuid.UUID) (*item, error) {
	return sv.repo.GetItem(ctx, id)
}

func (sv service) CreateItem(ctx context.Context, i *item) (*uuid.UUID, error) {
	return sv.repo.CreateItem(ctx, i)
}

func (sv service) ToggleItemComplete(ctx context.Context, id *uuid.UUID) (isDone, error) {
	item, err := sv.repo.GetItem(ctx, id)
	if err != nil {
		return false, err
	}

	item.IsDone = !item.IsDone
	newStatus := item.IsDone

	err = sv.repo.UpdateItem(ctx, id, item)
	if err != nil {
		return false, err
	}

	return newStatus, nil
}
func (sv service) UpdateItem(ctx context.Context, id *uuid.UUID, title string, description string) (*item, error) {
	item, err := sv.repo.GetItem(ctx, id)
	if err != nil {
		return nil, err
	}

	item.Title = title
	item.Description = description

	err = sv.repo.UpdateItem(ctx, id, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (sv service) DeleteItem(ctx context.Context, id *uuid.UUID) error {
	return sv.repo.DeleteItem(ctx, id)
}
