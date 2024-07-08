package grpc

import (
	"context"
	"errors"

	"github.com/dokyan1989/goit/app/grpc/pb"
	"github.com/dokyan1989/goit/internal/todo"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
)

type service struct {
	pb.UnimplementedTodoManagementServer
	svc *todo.Service
}

func New(db *pgxpool.Pool) *service {
	return &service{svc: todo.NewService(db)}
}

// CreateTodo
func (s *service) CreateTodo(ctx context.Context, request *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	params := todo.CreateTodoParams{
		Title:  request.Title,
		Status: request.Status,
	}
	id, err := s.svc.CreateTodo(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTodoResponse{
		Message: codes.OK.String(),
		Data:    &pb.CreateTodoResponse_Data{Id: id},
	}, nil
}

// UpdateTodo
func (s *service) UpdateTodo(ctx context.Context, request *pb.UpdateTodoRequest) (*pb.UpdateTodoResponse, error) {
	params := todo.UpdateTodoParams{
		Title:  request.Title,
		Status: request.Status,
	}
	err := s.svc.UpdateTodo(ctx, request.Id, params)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateTodoResponse{
		Message: codes.OK.String(),
		Data:    &pb.UpdateTodoResponse_Data{Id: request.Id},
	}, nil
}

// GetTodo
func (s *service) GetTodo(ctx context.Context, request *pb.GetTodoRequest) (*pb.GetTodoResponse, error) {
	params := todo.FindTodosParams{ID: &request.Id}
	todos, err := s.svc.FindTodos(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(todos) == 0 {
		return nil, errors.New("record not found")
	}

	pbTodo := pb.Todo{
		Id:     todos[0].ID,
		Title:  todos[0].Title,
		Status: todos[0].Status,
	}

	return &pb.GetTodoResponse{
		Message: codes.OK.String(),
		Data:    &pb.GetTodoResponse_Data{Todo: &pbTodo},
	}, nil
}

// ListTodos
func (s *service) ListTodos(ctx context.Context, request *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	params := todo.FindTodosParams{
		IDs: request.GetIds(),
		// Status: &request.GetStatus().GetValue(),
		Limit:  int(request.GetLimit()),
		Offset: int(request.GetOffset()),
	}
	todos, err := s.svc.FindTodos(ctx, params)
	if err != nil {
		return nil, err
	}

	pbTodos := make([]*pb.Todo, len(todos))
	for i, todo := range todos {
		pbTodos[i] = &pb.Todo{
			Id:     todo.ID,
			Title:  todo.Title,
			Status: todo.Status,
		}
	}

	return &pb.ListTodosResponse{
		Message: codes.OK.String(),
		Data:    &pb.ListTodosResponse_Data{Todos: pbTodos},
	}, nil
}

// DeleteTodo
func (s *service) DeleteTodo(ctx context.Context, request *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	err := s.svc.DeleteTodo(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.DeleteTodoResponse{
		Message: codes.OK.String(),
		Data:    &pb.DeleteTodoResponse_Data{Id: request.Id},
	}, nil
}
