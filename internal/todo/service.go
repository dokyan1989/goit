package todo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	store *postgres
}

func NewService(db *pgxpool.Pool) *Service {
	store := &postgres{db: db}

	return &Service{store: store}
}

func (s *Service) CreateTodo(ctx context.Context, params CreateTodoParams) (int64, error) {
	id, err := s.store.CreateTodo(ctx, params)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) UpdateTodo(ctx context.Context, id int64, params UpdateTodoParams) error {
	err := s.store.UpdateTodo(ctx, id, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) FindTodos(ctx context.Context, params FindTodosParams) ([]Todo, error) {
	todos, err := s.store.FindTodos(ctx, params)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (s *Service) DeleteTodo(ctx context.Context, id int64) error {
	err := s.store.DeleteTodo(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
