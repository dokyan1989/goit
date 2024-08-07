package helloweb

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	todostore "github.com/dokyan1989/goit/internal/todo"
	"github.com/dokyan1989/goit/misc/httpparser"
	"github.com/rs/zerolog"
)

/**
|-------------------------------------------------------------------------
| ViewTodo
|-----------------------------------------------------------------------*/

func (s *server) ViewTodo(w http.ResponseWriter, r *http.Request) *Page {
	logger := zerolog.Ctx(r.Context()).With().Str("api", "ViewTodo").Logger()

	var (
		id  int64
		err error
	)

	if err = httpparser.RouteParams(r, "id", &id); err != nil {
		err2 := fmt.Errorf("parse route params %q: %w", "id", err)
		return pageError(err2)
	}

	logger.Info().Int64("id", id).Send()

	res, err := viewTodo(r.Context(), s.todoStore, id)
	if err != nil {
		err2 := fmt.Errorf("get todo error: %w", err)
		return pageError(err2)
	}

	res.Title = "Todo Details"
	logger.Info().Str("response info", convertToJSON(res)).Send()

	return page("view", res)
}

type ViewTodoResponse struct {
	Title string
	Todo
}

func viewTodo(ctx context.Context, store TodoStore, id int64) (*ViewTodoResponse, error) {
	params := todostore.FindTodosParams{ID: id, Limit: 1}
	storeTodos, err := store.FindTodos(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(storeTodos) == 0 {
		return nil, errors.New("record not found")
	}

	st := storeTodos[0]
	todo := Todo{ID: st.ID, Title: st.Title, Status: st.Status}

	return &ViewTodoResponse{Todo: todo}, nil
}
