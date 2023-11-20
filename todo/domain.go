package todo

import (
	"time"

	"github.com/google/uuid"
)

type list struct {
	ID        uuid.UUID `redis:"id"`
	CreatedAt time.Time `redis:"createdAt"`
	Items     []*item   `redis:"createdAt"`
}

func newList() *list {
	return &list{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}
}

type item struct {
	ID          uuid.UUID `redis:"id"`
	CreatedAt   time.Time `redis:"createdAt"`
	Title       string    `redis:"title"`
	Description string    `redis:"description"`
	IsDone      isDone    `redis:"isDone"`
}

func newItem(title string, description string) *item {
	return &item{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		Title:       title,
		Description: description,
		IsDone:      false,
	}
}

type isDone bool
