package api

import "context"

const (
	MsgSuccess = "success"
)

/**
|-------------------------------------------------------------------------
| API
|-----------------------------------------------------------------------*/

// APIReponse ...
type APIReponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewAPIResponse(msg string, data any) *APIReponse {
	return &APIReponse{msg, data}
}

// APIError ...
type APIError struct {
	Err     error          `json:"error"`
	Details map[string]any `json:"details"`
}

func NewAPIError(err error) *APIError {
	return &APIError{Err: err}
}

func (err *APIError) Error() string {
	return err.Err.Error()
}

func (err *APIError) WithDetails(k string, v any) {
	if err.Details == nil {
		err.Details = make(map[string]any)
	}

	err.Details[k] = v
}

/**
|-------------------------------------------------------------------------
| TODO
|-----------------------------------------------------------------------*/

type Todo struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type TodoService interface {
	ListTodos(ctx context.Context, request *ListTodosRequest) (*ListTodosResponse, error)
	GetTodo(ctx context.Context, request *GetTodoRequest) (*GetTodoResponse, error)
	CreateTodo(ctx context.Context, request *CreateTodoRequest) (*CreateTodoResponse, error)
	UpdateTodo(ctx context.Context, request *UpdateTodoRequest) (*UpdateTodoResponse, error)
	DeleteTodo(ctx context.Context, request *DeleteTodoRequest) (*DeleteTodoResponse, error)
}

type CreateTodoRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

type CreateTodoResponse struct {
	ID int64 `json:"id"`
}

type UpdateTodoRequest struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type UpdateTodoResponse struct {
	ID int64 `json:"id"`
}

type GetTodoRequest struct {
	ID int64 `json:"id"`
}

type GetTodoResponse struct {
	Todo Todo `json:"todo"`
}

type ListTodosRequest struct {
	IDs    []int64 `schema:"ids"`
	Status string  `schema:"status"`
	Limit  int     `schema:"limit"`
	Offset int     `schema:"offset"`
}

type ListTodosResponse struct {
	Todos []Todo `json:"todos"`
}

type DeleteTodoRequest struct {
	ID int64 `json:"id"`
}

type DeleteTodoResponse struct {
	ID int64 `json:"id"`
}
