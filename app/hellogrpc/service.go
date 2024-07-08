package grpc

import (
	"context"
	"errors"

	"github.com/dokyan1989/goit/app/hellogrpc/pb"
	todostore "github.com/dokyan1989/goit/internal/todo"
	"google.golang.org/grpc/codes"
)

type TodoStore interface {
	CreateTodo(ctx context.Context, params todostore.CreateTodoParams) (int64, error)
	UpdateTodo(ctx context.Context, id int64, params todostore.UpdateTodoParams) error
	FindTodos(ctx context.Context, params todostore.FindTodosParams) ([]todostore.Todo, error)
	DeleteTodo(ctx context.Context, id int64) error
}

type todoSvc struct {
	pb.UnimplementedTodoManagementServer
	todoStore TodoStore
}

var _ pb.TodoManagementServer = (*todoSvc)(nil)

// CreateTodo
func (s *todoSvc) CreateTodo(ctx context.Context, request *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	params := todostore.CreateTodoParams{
		Title:  request.Title,
		Status: request.Status,
	}
	id, err := s.todoStore.CreateTodo(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTodoResponse{
		Message: codes.OK.String(),
		Data:    &pb.CreateTodoResponse_Data{Id: id},
	}, nil
}

// UpdateTodo
func (s *todoSvc) UpdateTodo(ctx context.Context, request *pb.UpdateTodoRequest) (*pb.UpdateTodoResponse, error) {
	params := todostore.UpdateTodoParams{
		Title:  request.Title,
		Status: request.Status,
	}
	err := s.todoStore.UpdateTodo(ctx, request.Id, params)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateTodoResponse{
		Message: codes.OK.String(),
		Data:    &pb.UpdateTodoResponse_Data{Id: request.Id},
	}, nil
}

// GetTodo
func (s *todoSvc) GetTodo(ctx context.Context, request *pb.GetTodoRequest) (*pb.GetTodoResponse, error) {
	params := todostore.FindTodosParams{ID: request.Id}
	storeTodos, err := s.todoStore.FindTodos(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(storeTodos) == 0 {
		return nil, errors.New("record not found")
	}

	todo := pb.Todo{Id: storeTodos[0].ID, Title: storeTodos[0].Title, Status: storeTodos[0].Status}

	return &pb.GetTodoResponse{
		Message: codes.OK.String(),
		Data:    &pb.GetTodoResponse_Data{Todo: &todo},
	}, nil
}

// ListTodos
func (s *todoSvc) ListTodos(ctx context.Context, request *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	params := todostore.FindTodosParams{
		IDs:    request.GetIds(),
		Status: request.GetStatus().GetValue(),
		Limit:  int(request.GetLimit()),
		Offset: int(request.GetOffset()),
	}
	storeTodos, err := s.todoStore.FindTodos(ctx, params)
	if err != nil {
		return nil, err
	}

	todos := make([]*pb.Todo, len(storeTodos))
	for i, st := range storeTodos {
		todos[i] = &pb.Todo{Id: st.ID, Title: st.Title, Status: st.Status}
	}

	return &pb.ListTodosResponse{
		Message: codes.OK.String(),
		Data:    &pb.ListTodosResponse_Data{Todos: todos},
	}, nil
}

// DeleteTodo
func (s *todoSvc) DeleteTodo(ctx context.Context, request *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	err := s.todoStore.DeleteTodo(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.DeleteTodoResponse{
		Message: codes.OK.String(),
		Data:    &pb.DeleteTodoResponse_Data{Id: request.Id},
	}, nil
}
