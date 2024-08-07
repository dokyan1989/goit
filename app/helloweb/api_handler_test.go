package helloweb

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	todostore "github.com/dokyan1989/goit/internal/todo"
	"github.com/dokyan1989/goit/internal/todo/todotest"
	"github.com/dokyan1989/goit/misc/t/httprequest"
	"github.com/go-chi/chi/middleware"
	"github.com/google/go-cmp/cmp"
)

// ref: https://semaphoreci.com/blog/table-driven-unit-tests-go
func TestHandler_CreateTodo(t *testing.T) {
	type Input struct {
		httpRequestOpts httprequest.Options
		mockCreatedID   int64
		mockErr         error
	}

	type Want struct {
		createTodoParams todostore.CreateTodoParams
		status           int
		response         string
	}

	tests := []struct {
		name  string
		input Input
		want  Want
	}{
		{
			name: "should return success",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Body: map[string]any{"title": "title", "status": TodoStatusUnknown},
				},
				mockCreatedID: 1,
			},
			want: Want{
				createTodoParams: todostore.CreateTodoParams{Title: "title", Status: TodoStatusUnknown},
				status:           http.StatusCreated,
				response:         `{"data":{"id":1},"message":"ok","request_id":"123"}`,
			},
		},
		{
			name: "should return error from store",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Body: map[string]any{"title": "title", "status": TodoStatusUnknown},
				},
				mockErr: errors.New("error from store"),
			},
			want: Want{
				createTodoParams: todostore.CreateTodoParams{Title: "title", Status: TodoStatusUnknown},
				status:           http.StatusInternalServerError,
				response:         `{"message":"failed to create the record","request_id":"123"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare
			var capture struct {
				mockParams todostore.CreateTodoParams
			}

			mockTodoStore := &todotest.Mock{
				CreateTodoFunc: func(ctx context.Context, params todostore.CreateTodoParams) (int64, error) {
					capture.mockParams = params
					return tt.input.mockCreatedID, tt.input.mockErr
				},
			}
			s := &server{todoStore: mockTodoStore}
			w := httptest.NewRecorder()
			r := httprequest.MustNewTest(t, http.MethodPost, "/api/v1/todos", tt.input.httpRequestOpts)
			r = r.WithContext(context.WithValue(r.Context(), middleware.RequestIDKey, "123"))

			// action
			makeAPIHandler(s.CreateTodo)(w, r)

			// assert resutls
			if diff := cmp.Diff(tt.want.createTodoParams, capture.mockParams); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
				return
			}

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("response status code mismatch (-want +got):\n-\t%d\n+\t%d",
					tt.want.status, res.StatusCode,
				)
			}

			gotBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("failed to read response body: %v", err)
				return
			}
			wantBody := []byte(tt.want.response)

			if diff := cmp.Diff(wantBody, gotBody, responseCmp()); diff != "" {
				t.Errorf("respone body mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandler_UpdateTodo(t *testing.T) {
	type Input struct {
		httpRequestOpts httprequest.Options
		mockErr         error
	}

	type Want struct {
		updateTodoID     int64
		updateTodoParams todostore.UpdateTodoParams
		status           int
		response         string
	}

	tests := []struct {
		name  string
		input Input
		want  Want
	}{
		{
			name: "should return success",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Params: []string{"id", "1"},
					Body:   map[string]any{"title": "title", "status": TodoStatusDoing},
				},
			},
			want: Want{
				updateTodoID:     1,
				updateTodoParams: todostore.UpdateTodoParams{Title: "title", Status: TodoStatusDoing},
				status:           http.StatusOK,
				response:         `{"data":{"id":1},"message":"ok","request_id":"123"}`,
			},
		},
		{
			name: "should return error parsing route params failed",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Params: []string{"id", "invalid"},
					Body:   map[string]any{"title": "title", "status": TodoStatusDoing},
				},
			},
			want: Want{
				status:   http.StatusBadRequest,
				response: `{"message":"failed to parse route params 'id'","request_id":"123"}`,
			},
		},
		{
			name: "should return error from store",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Params: []string{"id", "1"},
					Body:   map[string]any{"title": "title", "status": TodoStatusDoing},
				},
				mockErr: errors.New("error from store"),
			},
			want: Want{
				updateTodoID:     1,
				updateTodoParams: todostore.UpdateTodoParams{Title: "title", Status: TodoStatusDoing},
				status:           http.StatusInternalServerError,
				response:         `{"message":"failed to update the record","request_id":"123"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup mocked server
			var capture struct {
				mockID     int64
				mockParams todostore.UpdateTodoParams
			}
			mockTodoStore := &todotest.Mock{
				UpdateTodoFunc: func(ctx context.Context, id int64, params todostore.UpdateTodoParams) error {
					capture.mockID = id
					capture.mockParams = params
					return tt.input.mockErr
				},
			}
			s := &server{todoStore: mockTodoStore}

			w := httptest.NewRecorder()
			r := httprequest.MustNewTest(t, http.MethodPut, "/api/v1/todos/{id}", tt.input.httpRequestOpts)
			r = r.WithContext(context.WithValue(r.Context(), middleware.RequestIDKey, "123"))
			makeAPIHandler(s.UpdateTodo)(w, r)

			// assert resutls
			if capture.mockID != tt.want.updateTodoID {
				t.Errorf("result mismatch (-want:%d +got:%d):\n", tt.want.updateTodoID, capture.mockID)
				return
			}

			if diff := cmp.Diff(tt.want.updateTodoParams, capture.mockParams); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
				return
			}

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("response status code mismatch (-want +got):\n-\t%d\n+\t%d",
					tt.want.status, res.StatusCode,
				)
			}

			gotBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("failed to read response body: %v", err)
				return
			}
			wantBody := []byte(tt.want.response)

			if diff := cmp.Diff(wantBody, gotBody, responseCmp()); diff != "" {
				t.Errorf("respone body mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandler_GetTodo(t *testing.T) {
	type Input struct {
		httpRequestOpts httprequest.Options
		mockTodos       []todostore.Todo
		mockErr         error
	}

	type Want struct {
		findTodoParams todostore.FindTodosParams
		status         int
		response       string
	}

	tests := []struct {
		name  string
		input Input
		want  Want
	}{
		{
			name: "should return a todo successfully",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Params: []string{"id", "1"},
				},
				mockTodos: []todostore.Todo{
					{ID: 1, Title: "title", Status: TodoStatusUnknown},
				},
			},
			want: Want{
				findTodoParams: todostore.FindTodosParams{ID: 1},
				status:         http.StatusOK,
				response:       `{"data":{"todo":{"id":1,"status":"UNKNOWN","title":"title"}},"message":"ok","request_id":"123"}`,
			},
		},
		{
			name: "should return error parsing route params failed",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Params: []string{"id", "invalid"},
				},
			},
			want: Want{
				status:   http.StatusBadRequest,
				response: `{"message":"failed to parse route params 'id'","request_id":"123"}`,
			},
		},
		{
			name: "should return error from store",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Params: []string{"id", "1"},
				},
				mockErr: errors.New("error from store"),
			},
			want: Want{
				findTodoParams: todostore.FindTodosParams{ID: 1},
				status:         http.StatusInternalServerError,
				response:       `{"message":"failed to get the record","request_id":"123"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup mocked server
			var capture struct {
				mockParams todostore.FindTodosParams
			}

			mockTodoStore := &todotest.Mock{
				FindTodosFunc: func(ctx context.Context, params todostore.FindTodosParams) ([]todostore.Todo, error) {
					capture.mockParams = params
					return tt.input.mockTodos, tt.input.mockErr
				},
			}
			s := &server{todoStore: mockTodoStore}

			// action
			w := httptest.NewRecorder()
			r := httprequest.MustNewTest(t, http.MethodGet, "/api/v1/todos/{id}", tt.input.httpRequestOpts)
			r = r.WithContext(context.WithValue(r.Context(), middleware.RequestIDKey, "123"))
			makeAPIHandler(s.GetTodo)(w, r)

			// assert results
			if diff := cmp.Diff(tt.want.findTodoParams, capture.mockParams); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
				return
			}

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("response status code mismatch (-want +got):\n-\t%d\n+\t%d",
					tt.want.status, res.StatusCode,
				)
			}

			gotBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("failed to read response body: %v", err)
				return
			}
			wantBody := []byte(tt.want.response)

			if diff := cmp.Diff(wantBody, gotBody, responseCmp()); diff != "" {
				t.Errorf("respone body mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandler_ListTodo(t *testing.T) {
	type Input struct {
		httpRequestOpts httprequest.Options
		mockTodos       []todostore.Todo
		mockErr         error
	}

	type Want struct {
		findTodosParams todostore.FindTodosParams
		status          int
		response        string
	}

	tests := []struct {
		name  string
		input Input
		want  Want
	}{
		{
			name: "should return a list of todos successfully",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Query: url.Values{"ids": []string{"1", "2"}},
				},
				mockTodos: []todostore.Todo{
					{ID: 1, Title: "title_1", Status: TodoStatusUnknown},
					{ID: 2, Title: "title_2", Status: TodoStatusUnknown},
				},
			},
			want: Want{
				findTodosParams: todostore.FindTodosParams{IDs: []int64{1, 2}, Status: ""},
				status:          http.StatusOK,
				response:        `{"data":{"todos":[{"id":1,"status":"UNKNOWN","title":"title_1"},{"id":2,"status":"UNKNOWN","title":"title_2"}]},"message":"ok","request_id":"123"}`,
			},
		},
		{
			name: "should return error parsing query params failed",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Query: url.Values{"idss": []string{"1", "2"}},
				},
			},
			want: Want{
				status:   http.StatusBadRequest,
				response: `{"message":"invalid: query parameters is empty","request_id":"123"}`,
			},
		},
		{
			name: "should return error from store",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Query: url.Values{"ids": []string{"1", "2"}},
				},
				mockErr: errors.New("error from store"),
			},
			want: Want{
				findTodosParams: todostore.FindTodosParams{IDs: []int64{1, 2}, Status: ""},
				status:          http.StatusInternalServerError,
				response:        `{"message":"failed to list the records","request_id":"123"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup mocked server
			var capture struct {
				mockParams todostore.FindTodosParams
			}

			mockTodoStore := &todotest.Mock{
				FindTodosFunc: func(ctx context.Context, params todostore.FindTodosParams) ([]todostore.Todo, error) {
					capture.mockParams = params
					return tt.input.mockTodos, tt.input.mockErr
				},
			}
			s := &server{todoStore: mockTodoStore}

			// action
			w := httptest.NewRecorder()
			r := httprequest.MustNewTest(t, http.MethodGet, "/api/v1/todos", tt.input.httpRequestOpts)
			r = r.WithContext(context.WithValue(r.Context(), middleware.RequestIDKey, "123"))
			makeAPIHandler(s.ListTodos)(w, r)

			// assert results
			if diff := cmp.Diff(tt.want.findTodosParams, capture.mockParams); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
				return
			}

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("response status code mismatch (-want +got):\n-\t%d\n+\t%d",
					tt.want.status, res.StatusCode,
				)
			}

			gotBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("failed to read response body: %v", err)
				return
			}
			wantBody := []byte(tt.want.response)

			if diff := cmp.Diff(wantBody, gotBody, responseCmp()); diff != "" {
				t.Errorf("respone body mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandler_DeleteTodo(t *testing.T) {
	type Input struct {
		httpRequestOpts httprequest.Options
		mockErr         error
	}

	type Want struct {
		deleteTodoID int64
		status       int
		response     string
	}

	tests := []struct {
		name  string
		input Input
		want  Want
	}{
		{
			name: "should return a todo record",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Params: []string{"id", "1"},
				},
			},
			want: Want{
				deleteTodoID: 1,
				status:       http.StatusOK,
				response:     `{"data":{"id":1},"message":"ok","request_id":"123"}`,
			},
		},
		{
			name: "should return error parsing route params failed",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Params: []string{"id", "invalid"},
				},
			},
			want: Want{
				status:   http.StatusBadRequest,
				response: `{"message":"failed to parse route params 'id'","request_id":"123"}`,
			},
		},
		{
			name: "should return error from store",
			input: Input{
				httpRequestOpts: httprequest.Options{
					Params: []string{"id", "1"},
				},
				mockErr: errors.New("error from store"),
			},
			want: Want{
				deleteTodoID: 1,
				status:       http.StatusInternalServerError,
				response:     `{"message":"failed to delete the record","request_id":"123"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup mocked server
			var capture struct {
				mockID int64
			}

			mockTodoStore := &todotest.Mock{
				DeleteTodoFunc: func(ctx context.Context, id int64) error {
					capture.mockID = id
					return tt.input.mockErr
				},
			}
			s := &server{todoStore: mockTodoStore}

			// action
			w := httptest.NewRecorder()
			r := httprequest.MustNewTest(t, http.MethodDelete, "/api/v1/todos/{id}", tt.input.httpRequestOpts)
			r = r.WithContext(context.WithValue(r.Context(), middleware.RequestIDKey, "123"))
			makeAPIHandler(s.DeleteTodo)(w, r)

			// assert results
			if diff := cmp.Diff(tt.want.deleteTodoID, capture.mockID); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
				return
			}

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("response status code mismatch (-want +got):\n-\t%d\n+\t%d",
					tt.want.status, res.StatusCode,
				)
			}

			gotBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("failed to read response body: %v", err)
				return
			}
			wantBody := []byte(tt.want.response)

			if diff := cmp.Diff(wantBody, gotBody, responseCmp()); diff != "" {
				t.Errorf("respone body mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
