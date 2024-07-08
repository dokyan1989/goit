package helloapi

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
| createTodo
|-----------------------------------------------------------------------*/

func (s *server) CreateTodo(w http.ResponseWriter, r *http.Request) (*APIResponse, error) {
	logger := zerolog.Ctx(r.Context())

	var body struct {
		Title  string `json:"title"`
		Status string `json:"status"`
	}

	status, err := httpparser.RequestBody(w, r, &body)
	if err != nil {
		err2 := fmt.Errorf("[Create Todo] parse request body: %w", err)
		return nil, apiError(err2, "failed to parse request body").WithStatus(status)
	}

	logger.Info().Str("[Create Todo] body", convertToJSON(body)).Send()

	p := todostore.CreateTodoParams{
		Title:  body.Title,
		Status: body.Status,
	}
	res, err := createTodo(r.Context(), s.todoStore, p)
	if err != nil {
		err2 := fmt.Errorf("[Create Todo] store error: %w", err)
		return nil, apiError(err2, "failed to create the record")
	}

	apiRes := apiResponse(MsgSuccess, res).WithStatus(http.StatusCreated)

	logger.Info().Str("[Create Todo] response", convertToJSON(apiRes)).Send()

	return apiRes, nil
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
| updateTodo
|-----------------------------------------------------------------------*/

func (s *server) UpdateTodo(w http.ResponseWriter, r *http.Request) (*APIResponse, error) {
	logger := zerolog.Ctx(r.Context())

	var id int64
	var body struct {
		Title  string `json:"title"`
		Status string `json:"status"`
	}

	status, err := httpparser.RequestBody(w, r, &body)
	if err != nil {
		err2 := fmt.Errorf("[Update Todo] parse request body: %w", err)
		return nil, apiError(err2, "failed to parse request body").WithStatus(status)
	}

	if err = httpparser.RouteParams(r, "id", &id); err != nil {
		err2 := fmt.Errorf("[Update Todo] parse route params %q: %w", "id", err)
		return nil, apiError(err2, "failed to parse route params 'id'").WithStatus(http.StatusBadRequest)
	}

	logger.Info().Str("[Update Todo] body", convertToJSON(body)).Int64("id", id).Send()

	p := todostore.UpdateTodoParams{
		Title:  body.Title,
		Status: body.Status,
	}
	res, err := updateTodo(r.Context(), s.todoStore, id, p)
	if err != nil {
		err2 := fmt.Errorf("[Update Todo] store error: %w", err)
		return nil, apiError(err2, "failed to update the record")
	}

	apiRes := apiResponse(MsgSuccess, res)

	logger.Info().Str("[Update Todo] response", convertToJSON(apiRes)).Send()

	return apiRes, nil
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
| getTodo
|-----------------------------------------------------------------------*/

func (s *server) GetTodo(w http.ResponseWriter, r *http.Request) (*APIResponse, error) {
	logger := zerolog.Ctx(r.Context())

	var (
		id  int64
		err error
	)

	if err = httpparser.RouteParams(r, "id", &id); err != nil {
		err2 := fmt.Errorf("[Get Todo] parse route params %q: %w", "id", err)
		return nil, apiError(err2, "failed to parse route params 'id'").WithStatus(http.StatusBadRequest)
	}

	logger.Info().Int64("[Get Todo] id", id).Send()

	res, err := getTodo(r.Context(), s.todoStore, id)
	if err != nil {
		err2 := fmt.Errorf("[Get Todo] store error: %w", err)
		return nil, apiError(err2, "failed to get the record")
	}

	apiRes := apiResponse(MsgSuccess, res)

	logger.Info().Str("[Get Todo] response", convertToJSON(apiRes)).Send()

	return apiRes, nil
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
| listTodos
|-----------------------------------------------------------------------*/

func (s *server) ListTodos(w http.ResponseWriter, r *http.Request) (*APIResponse, error) {
	logger := zerolog.Ctx(r.Context())

	var q struct {
		IDs    []int64 `json:"ids"`
		Status string  `json:"status"`
		Limit  int     `json:"limit"`
		Offset int     `json:"offset"`
	}

	err := httpparser.RequestQuery(r, &q)
	if err != nil {
		err2 := fmt.Errorf("[List Todos] parse request query: %w", err)
		return nil, apiError(err2, "failed to parse request query").WithStatus(http.StatusBadRequest)
	}

	logger.Info().Str("[List Todos] query", convertToJSON(q)).Send()

	// https://freshman.tech/snippets/go/check-empty-struct/
	if reflect.ValueOf(q).IsZero() {
		err2 := errors.New("invalid: query parameters is empty")
		return nil, apiError(err2, err2.Error()).WithStatus(http.StatusBadRequest)
	}

	p := todostore.FindTodosParams{
		IDs:    q.IDs,
		Status: q.Status,
		Limit:  q.Limit,
		Offset: q.Offset,
	}
	res, err := listTodos(r.Context(), s.todoStore, p)
	if err != nil {
		err2 := fmt.Errorf("[List Todos] store error: %w", err)
		return nil, apiError(err2, "failed to list the records")
	}

	apiRes := apiResponse(MsgSuccess, res)

	logger.Info().Str("[List Todos] response", convertToJSON(apiRes)).Send()

	return apiRes, nil
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
| deleteTodo
|-----------------------------------------------------------------------*/

func (s *server) DeleteTodo(w http.ResponseWriter, r *http.Request) (*APIResponse, error) {
	logger := zerolog.Ctx(r.Context())

	var (
		id  int64
		err error
	)

	if err = httpparser.RouteParams(r, "id", &id); err != nil {
		err2 := fmt.Errorf("[Delete Todo] parse route params %q: %w", "id", err)
		return nil, apiError(err2, "failed to parse route params 'id'").WithStatus(http.StatusBadRequest)
	}

	logger.Info().Int64("[Delete Todo] id", id).Send()

	res, err := deleteTodo(r.Context(), s.todoStore, id)
	if err != nil {
		err2 := fmt.Errorf("[Delete Todo] store error: %w", err)
		return nil, apiError(err2, "failed to delete the record")
	}

	apiRes := apiResponse(MsgSuccess, res)

	logger.Info().Str("[Delete Todo] response", convertToJSON(apiRes)).Send()

	return apiRes, nil
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
