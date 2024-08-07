package helloweb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	todostore "github.com/dokyan1989/goit/internal/todo"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

const (
	MsgOK = "ok"
)

/**
|-------------------------------------------------------------------------
| API
|-----------------------------------------------------------------------*/

// API ...
type API struct {
	Status    int
	Message   string
	Data      any
	Err       error
	RequestID string
	Extra     map[string]any
}

func apiOK(data any) *API {
	return &API{Status: http.StatusOK, Message: MsgOK, Data: data}
}

func apiError(err error) *API {
	return &API{Status: http.StatusInternalServerError, Err: err}
}

func (api *API) WithStatus(status int) *API {
	if api != nil {
		api.Status = status
	}
	return api
}

func (api *API) WithMessage(msg string) *API {
	if api != nil {
		api.Message = msg
	}
	return api
}

func (api *API) WithRequestID(id string) *API {
	if api != nil {
		api.RequestID = id
	}
	return api
}

func (api *API) WithExtra(extra map[string]any) *API {
	if api != nil {
		api.Extra = extra
	}
	return api
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if api == nil {
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(api.Status)

	type response struct {
		Message   string         `json:"message"`
		Extra     map[string]any `json:"extra,omitempty"`
		Data      any            `json:"data,omitempty"`
		RequestID string         `json:"request_id"`
	}

	// write error
	if api.Err != nil {
		l := zerolog.Ctx(r.Context())

		rawQuery := r.URL.RawQuery
		if rawQuery != "" {
			rawQuery = fmt.Sprintf("?%s", rawQuery)
		}
		l.Err(api.Err).Msg(fmt.Sprintf("error on API:[%s] %s%s", r.Method, r.URL.Path, rawQuery))

		msg := api.Message
		if msg == "" {
			msg = api.Err.Error()
		}

		resp := response{
			Message:   msg,
			Extra:     api.Extra,
			RequestID: api.RequestID,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// write ok
	resp := response{
		Message:   api.Message,
		Data:      api.Data,
		RequestID: api.RequestID,
	}
	json.NewEncoder(w).Encode(resp)
}

type RequestIDSetter interface {
	SetRequestID(id string)
}

type APIHandler func(w http.ResponseWriter, r *http.Request) *API

func makeAPIHandler(h APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())

		// inject request_id to the logger
		l := zerolog.Ctx(r.Context())
		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Logger().With().Str("request_id", reqID)
		})

		api := h(w, r).WithRequestID(reqID)
		api.ServeHTTP(w, r)
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
