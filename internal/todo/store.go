package todo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

type CreateTodoParams struct {
	Title  string
	Status string
}

func (s *postgres) CreateTodo(ctx context.Context, params CreateTodoParams) (int64, error) {
	var id int64

	row := s.db.QueryRow(ctx, "INSERT INTO todo (title, status) VALUES ($1, $2) RETURNING id", params.Title, params.Status)
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("no row inserted")
		}
		return 0, err
	}

	return id, nil
}

type UpdateTodoParams struct {
	Title  string
	Status string
}

func (s *postgres) UpdateTodo(ctx context.Context, id int64, params UpdateTodoParams) error {
	tag, err := s.db.Exec(ctx, "UPDATE todo SET title = $1, status = $2 WHERE id = $3", params.Title, params.Status, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no row updated")
	}

	return nil
}

type FindTodosParams struct {
	// optional
	ID     *int64
	IDs    []int64
	Title  *string
	Status *string

	// required
	Limit  int
	Offset int
}

func (s *postgres) FindTodos(ctx context.Context, params FindTodosParams) ([]Todo, error) {
	var todos []Todo

	rows, err := s.db.Query(ctx, "select * from todo")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (s *postgres) DeleteTodo(ctx context.Context, id int64) error {
	tag, err := s.db.Exec(ctx, "DELETE FROM todo WHERE id = $1", id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no row deleted")
	}

	return nil
}
