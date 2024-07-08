package todotest

import (
	"context"

	"github.com/dokyan1989/goit/internal/todo"
)

//go:generate moq -out store.go -rm -with-resets . TodoStore:Mock
type TodoStore interface {
	CreateTodo(ctx context.Context, params todo.CreateTodoParams) (int64, error)
	UpdateTodo(ctx context.Context, id int64, params todo.UpdateTodoParams) error
	FindTodos(ctx context.Context, params todo.FindTodosParams) ([]todo.Todo, error)
	DeleteTodo(ctx context.Context, id int64) error
}
