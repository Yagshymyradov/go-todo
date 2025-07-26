package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Yagshymyradov/go-todo/internal/todo/repository"
)

var (
	ErrEmptyTitle = errors.New("title cannot be empty")
	ErrNotFound = errors.New("todo not found")
)

type TodoService struct {
	repo repository.Repository
}

func New(repo repository.Repository) *TodoService {
	return &TodoService{repo: repo}
}

type Todo struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *TodoService) Create(ctx context.Context, title string) (int64, error) {
	// Validate title
	title = strings.TrimSpace(title)
	if title == "" {
		return 0, ErrEmptyTitle
	}

	return s.repo.Create(ctx, title)
}

func (s *TodoService) List(ctx context.Context) ([]Todo, error) {
	todos, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Convert repository todos to service todos
	result := make([]Todo, len(todos))
	for i, t := range todos {
		result[i] = Todo{
			ID: t.ID,
			Title: t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		}
	}

	return result, nil
}

func (s *TodoService) Toggle(ctx context.Context, id int64) error {
	todos, err := s.repo.List(ctx)
	if err != nil {
		return err
	}

	// Find todo and get current completed status
	var found bool
	var currentStatus bool
	for _, t := range todos {
		if t.ID == id {
			found = true
			currentStatus = t.Completed
			break
		}
	}

	if !found {
		return ErrNotFound
	}

	// Toggle completed status
	return s.repo.Toggle(ctx, id, !currentStatus)
}

func (s *TodoService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}