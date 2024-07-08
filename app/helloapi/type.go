package helloapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	todostore "github.com/dokyan1989/goit/internal/todo"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

const (
	MsgSuccess = "success"
)

/**
|-------------------------------------------------------------------------
| API
|-----------------------------------------------------------------------*/

// APIResponse ...
type APIResponse struct {
	status    int
	Message   string `json:"message"`
	Data      any    `json:"data"`
	RequestID string `json:"request_id"`
}

func apiResponse(msg string, data any) *APIResponse {
	return &APIResponse{http.StatusOK, msg, data, ""}
}

func (res *APIResponse) WithStatus(status int) *APIResponse {
	if res != nil {
		res.status = status
	}
	return res
}

func (res *APIResponse) WithRequestID(reqId string) *APIResponse {
	if res != nil {
		res.RequestID = reqId
	}
	return res
}

func (res *APIResponse) Status() int {
	if res != nil {
		return res.status
	}
	return http.StatusOK
}

// APIError ...
type APIError struct {
	ErrorMessage string         `json:"error_message"`
	Details      map[string]any `json:"details,omitempty"`
	RequestID    string         `json:"request_id"`

	// unexported fileds
	status int
	err    error // root error
}

func apiError(err error, msg string) *APIError {
	return &APIError{
		ErrorMessage: msg,
		Details:      nil,
		RequestID:    "",
		err:          err,
		status:       http.StatusInternalServerError,
	}
}

func (err *APIError) WithStatus(status int) *APIError {
	if err != nil {
		err.status = status
	}
	return err
}

func (err *APIError) WithDetails(k string, v any) {
	if err.Details == nil {
		err.Details = make(map[string]any)
	}
	err.Details[k] = v
}

func (err *APIError) WithRequestID(id string) *APIError {
	if err != nil {
		err.RequestID = id
	}
	return err
}

func (err *APIError) Error() string {
	if err != nil {
		return err.err.Error()
	}
	return ""
}

func (err *APIError) Status() int {
	if err != nil {
		return err.status
	}
	return http.StatusInternalServerError
}

// APIHandler ...
type APIHandler func(w http.ResponseWriter, r *http.Request) (*APIResponse, error)

func makeAPIHandler(f APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())

		// inject request_id to the logger
		l := zerolog.Ctx(r.Context())
		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Logger().With().Str("request_id", reqID)
		})

		res, err := f(w, r)
		if err == nil {
			writeJSON(w, res.Status(), res.WithRequestID(reqID))
			return
		}

		// create APIError if necessary
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			apiErr = apiError(err, err.Error())
		}

		l.Info().Err(apiErr.err).Msg(fmt.Sprintf("Error on [%s %s?%s]", r.Method, r.URL.Path, r.URL.RawQuery))

		// https://stackoverflow.com/questions/68792696/why-error-messages-shouldnt-end-with-a-punctuation-mark-in-go
		// https://jeremybytes.blogspot.com/2021/01/go-golang-error-handling-different.html
		writeError(w, apiErr.Status(), apiErr.WithRequestID(reqID))
	}
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

const (
	TodoStatusUnknown = "UNKNOWN"
	TodoStatusNew     = "NEW"
	TodoStatusDoing   = "DOING"
	TodoStatusPending = "PENDING"
	TodoStatusDone    = "DONE"
)

type TodoStore interface {
	CreateTodo(ctx context.Context, params todostore.CreateTodoParams) (int64, error)
	UpdateTodo(ctx context.Context, id int64, params todostore.UpdateTodoParams) error
	FindTodos(ctx context.Context, params todostore.FindTodosParams) ([]todostore.Todo, error)
	DeleteTodo(ctx context.Context, id int64) error
}
