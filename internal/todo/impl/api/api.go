package api

import (
	"context"
	"errors"

	"github.com/dokyan1989/goit/app/api"
	"github.com/dokyan1989/goit/internal/todo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type service struct {
	svc *todo.Service
}

func New(db *pgxpool.Pool) *service {
	return &service{svc: todo.NewService(db)}
}

func (s *service) CreateTodo(ctx context.Context, request *api.CreateTodoRequest) (*api.CreateTodoResponse, error) {
	params := todo.CreateTodoParams{
		Title:  request.Title,
		Status: request.Status,
	}
	id, err := s.svc.CreateTodo(ctx, params)
	if err != nil {
		return nil, err
	}

	return &api.CreateTodoResponse{ID: id}, nil
}

func (s *service) UpdateTodo(ctx context.Context, request *api.UpdateTodoRequest) (*api.UpdateTodoResponse, error) {
	params := todo.UpdateTodoParams{
		Title:  request.Title,
		Status: request.Status,
	}
	err := s.svc.UpdateTodo(ctx, request.ID, params)
	if err != nil {
		return nil, err
	}

	return &api.UpdateTodoResponse{ID: request.ID}, nil
}

func (s *service) GetTodo(ctx context.Context, request *api.GetTodoRequest) (*api.GetTodoResponse, error) {
	params := todo.FindTodosParams{ID: &request.ID}
	todos, err := s.svc.FindTodos(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(todos) == 0 {
		return nil, errors.New("record not found")
	}

	apiTodo := api.Todo{
		ID:     todos[0].ID,
		Title:  todos[0].Title,
		Status: todos[0].Status,
	}

	return &api.GetTodoResponse{Todo: apiTodo}, nil
}

func (s *service) ListTodos(ctx context.Context, request *api.ListTodosRequest) (*api.ListTodosResponse, error) {
	params := todo.FindTodosParams{
		IDs:    request.IDs,
		Status: &request.Status,
		Limit:  request.Limit,
		Offset: request.Offset,
	}
	todos, err := s.svc.FindTodos(ctx, params)
	if err != nil {
		return nil, err
	}

	apiTodos := make([]api.Todo, len(todos))
	for i, todo := range todos {
		apiTodos[i] = api.Todo{
			ID:     todo.ID,
			Title:  todo.Title,
			Status: todo.Status,
		}
	}

	return &api.ListTodosResponse{Todos: apiTodos}, nil
}

func (s *service) DeleteTodo(ctx context.Context, request *api.DeleteTodoRequest) (*api.DeleteTodoResponse, error) {
	err := s.svc.DeleteTodo(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return &api.DeleteTodoResponse{ID: request.ID}, nil
}
