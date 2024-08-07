package helloweb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	todostore "github.com/dokyan1989/goit/internal/todo"
	"github.com/dokyan1989/goit/misc/httpparser"
	"github.com/rs/zerolog"
)

/**
|-------------------------------------------------------------------------
| CreateTodo
|-----------------------------------------------------------------------*/

func (s *server) CreateTodo(w http.ResponseWriter, r *http.Request) *API {
	logger := zerolog.Ctx(r.Context()).With().Str("api", "CreateTodo").Logger()

	var request struct {
		Title  string `json:"title"`
		Status string `json:"status"`
	}

	status, err := httpparser.RequestBody(w, r, &request)
	if err != nil {
		err2 := fmt.Errorf("parse request body: %w", err)
		return apiError(err2).WithMessage("failed to parse request body").WithStatus(status)
	}

	logger.Info().Str("request info", convertToJSON(request)).Send()

	p := todostore.CreateTodoParams{
		Title:  request.Title,
		Status: request.Status,
	}
	res, err := createTodo(r.Context(), s.todoStore, p)
	if err != nil {
		err2 := fmt.Errorf("create todo error: %w", err)
		return apiError(err2).WithMessage("failed to create the record")
	}

	logger.Info().Str("response info", convertToJSON(res)).Send()

	return apiOK(res).WithStatus(http.StatusCreated)
}

type CreateTodoResponse struct {
	ID int64 `json:"id"`
}

func createTodo(ctx context.Context, store TodoStore, p todostore.CreateTodoParams) (*CreateTodoResponse, error) {
	id, err := store.CreateTodo(ctx, p)
	if err != nil {
		return nil, err
	}

	return &CreateTodoResponse{ID: id}, nil
}

/**
|-------------------------------------------------------------------------
| UpdateTodo
|-----------------------------------------------------------------------*/

func (s *server) UpdateTodo(w http.ResponseWriter, r *http.Request) *API {
	logger := zerolog.Ctx(r.Context()).With().Str("api", "UpdateTodo").Logger()

	var request struct {
		ID     int64  `json:"id"`
		Title  string `json:"title"`
		Status string `json:"status"`
	}

	status, err := httpparser.RequestBody(w, r, &request)
	if err != nil {
		err2 := fmt.Errorf("parse request body: %w", err)
		return apiError(err2).WithMessage("failed to parse request body").WithStatus(status)
	}

	if err = httpparser.RouteParams(r, "id", &request.ID); err != nil {
		err2 := fmt.Errorf("parse route params %q: %w", "id", err)
		return apiError(err2).WithMessage("failed to parse route params 'id'").WithStatus(http.StatusBadRequest)
	}

	logger.Info().Str("request info", convertToJSON(request)).Send()

	p := todostore.UpdateTodoParams{
		Title:  request.Title,
		Status: request.Status,
	}
	res, err := updateTodo(r.Context(), s.todoStore, request.ID, p)
	if err != nil {
		err2 := fmt.Errorf("update todo error: %w", err)
		return apiError(err2).WithMessage("failed to update the record")
	}

	logger.Info().Str("response info", convertToJSON(res)).Send()

	return apiOK(res)
}

type UpdateTodoResponse struct {
	ID int64 `json:"id"`
}

func updateTodo(ctx context.Context, store TodoStore, id int64, p todostore.UpdateTodoParams) (*UpdateTodoResponse, error) {
	err := store.UpdateTodo(ctx, id, p)
	if err != nil {
		return nil, err
	}

	return &UpdateTodoResponse{ID: id}, nil
}

/**
|-------------------------------------------------------------------------
| GetTodo
|-----------------------------------------------------------------------*/

func (s *server) GetTodo(w http.ResponseWriter, r *http.Request) *API {
	logger := zerolog.Ctx(r.Context()).With().Str("api", "GetTodo").Logger()

	var (
		id  int64
		err error
	)

	if err = httpparser.RouteParams(r, "id", &id); err != nil {
		err2 := fmt.Errorf("parse route params %q: %w", "id", err)
		return apiError(err2).WithMessage("failed to parse route params 'id'").WithStatus(http.StatusBadRequest)
	}

	logger.Info().Int64("id", id).Send()

	res, err := getTodo(r.Context(), s.todoStore, id)
	if err != nil {
		err2 := fmt.Errorf("get todo error: %w", err)
		return apiError(err2).WithMessage("failed to get the record")
	}

	logger.Info().Str("response info", convertToJSON(res)).Send()

	return apiOK(res)
}

type GetTodoResponse struct {
	Todo Todo `json:"todo"`
}

func getTodo(ctx context.Context, store TodoStore, id int64) (*GetTodoResponse, error) {
	params := todostore.FindTodosParams{ID: id}
	storeTodos, err := store.FindTodos(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(storeTodos) == 0 {
		return nil, errors.New("record not found")
	}

	st := storeTodos[0]
	todo := Todo{ID: st.ID, Title: st.Title, Status: st.Status}

	return &GetTodoResponse{Todo: todo}, nil
}

/**
|-------------------------------------------------------------------------
| ListTodos
|-----------------------------------------------------------------------*/

func (s *server) ListTodos(w http.ResponseWriter, r *http.Request) *API {
	logger := zerolog.Ctx(r.Context()).With().Str("api", "ListTodos").Logger()

	var query struct {
		IDs    []int64 `json:"ids"`
		Status string  `json:"status"`
		Limit  int     `json:"limit"`
		Offset int     `json:"offset"`
	}

	err := httpparser.RequestQuery(r, &query)
	if err != nil {
		err2 := fmt.Errorf("parse request query: %w", err)
		return apiError(err2).WithMessage("failed to parse request query").WithStatus(http.StatusBadRequest)
	}

	logger.Info().Str("query", convertToJSON(query)).Send()

	// https://freshman.tech/snippets/go/check-empty-struct/
	if reflect.ValueOf(query).IsZero() {
		err2 := errors.New("invalid: query parameters is empty")
		return apiError(err2).WithMessage(err2.Error()).WithStatus(http.StatusBadRequest)
	}

	p := todostore.FindTodosParams{
		IDs:    query.IDs,
		Status: query.Status,
		Limit:  query.Limit,
		Offset: query.Offset,
	}
	res, err := listTodos(r.Context(), s.todoStore, p)
	if err != nil {
		err2 := fmt.Errorf("list todos error: %w", err)
		return apiError(err2).WithMessage("failed to list the records")
	}

	logger.Info().Str("response info", convertToJSON(res)).Send()

	return apiOK(res)
}

type ListTodosResponse struct {
	Todos []Todo `json:"todos"`
}

func listTodos(ctx context.Context, store TodoStore, p todostore.FindTodosParams) (*ListTodosResponse, error) {
	storeTodo, err := store.FindTodos(ctx, p)
	if err != nil {
		return nil, err
	}

	todos := make([]Todo, len(storeTodo))
	for i, st := range storeTodo {
		todos[i] = Todo{ID: st.ID, Title: st.Title, Status: st.Status}
	}

	return &ListTodosResponse{Todos: todos}, nil
}

/**
|-------------------------------------------------------------------------
| DeleteTodo
|-----------------------------------------------------------------------*/

func (s *server) DeleteTodo(w http.ResponseWriter, r *http.Request) *API {
	logger := zerolog.Ctx(r.Context()).With().Str("api", "DeleteTodo").Logger()

	var (
		id  int64
		err error
	)

	if err = httpparser.RouteParams(r, "id", &id); err != nil {
		err2 := fmt.Errorf("parse route params %q: %w", "id", err)
		return apiError(err2).WithMessage("failed to parse route params 'id'").WithStatus(http.StatusBadRequest)
	}

	logger.Info().Int64("id", id).Send()

	res, err := deleteTodo(r.Context(), s.todoStore, id)
	if err != nil {
		err2 := fmt.Errorf("delete todo error: %w", err)
		return apiError(err2).WithMessage("failed to delete the record")
	}

	logger.Info().Str("response info", convertToJSON(res)).Send()

	return apiOK(res)
}

type DeleteTodoResponse struct {
	ID int64 `json:"id"`
}

func deleteTodo(ctx context.Context, store TodoStore, id int64) (*DeleteTodoResponse, error) {
	err := store.DeleteTodo(ctx, id)
	if err != nil {
		return nil, err
	}

	return &DeleteTodoResponse{ID: id}, nil
}
