package helloweb

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	todostore "github.com/dokyan1989/goit/internal/todo"
	"github.com/dokyan1989/goit/internal/todo/todotest"
	"github.com/dokyan1989/goit/misc/t/httprequest"
	"github.com/go-chi/chi/middleware"
	"github.com/google/go-cmp/cmp"
	"github.com/sebdah/goldie/v2"
)

func TestHandler_ViewTodo(t *testing.T) {
	type Input struct {
		httpRequestOpts httprequest.Options
		mockTodos       []todostore.Todo
		mockErr         error
	}

	type Want struct {
		findTodoParams todostore.FindTodosParams
		status         int
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
				findTodoParams: todostore.FindTodosParams{ID: 1, Limit: 1},
				status:         http.StatusOK,
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
				status: http.StatusInternalServerError,
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
				findTodoParams: todostore.FindTodosParams{ID: 1, Limit: 1},
				status:         http.StatusInternalServerError,
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
			r := httprequest.MustNewTest(t, http.MethodGet, "/v1/todos/{id}", tt.input.httpRequestOpts)
			r = r.WithContext(context.WithValue(r.Context(), middleware.RequestIDKey, "123"))
			makePageHandler(s.ViewTodo)(w, r)

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

			g := goldie.New(t)
			g.Assert(t, t.Name(), w.Body.Bytes())
		})
	}
}
