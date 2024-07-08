package todo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

type CreateTodoParams struct {
	Title  string
	Status string
}

func (s *Store) CreateTodo(ctx context.Context, params CreateTodoParams) (int64, error) {
	id, err := createTodo(ctx, s.db, params)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func createTodo(ctx context.Context, db *pgxpool.Pool, params CreateTodoParams) (int64, error) {
	var id int64

	row := db.QueryRow(ctx, "INSERT INTO todo (title, status) VALUES ($1, $2) RETURNING id", params.Title, params.Status)
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

func (s *Store) UpdateTodo(ctx context.Context, id int64, params UpdateTodoParams) error {
	err := updateTodo(ctx, s.db, id, params)
	if err != nil {
		return err
	}

	return nil
}

func updateTodo(ctx context.Context, db *pgxpool.Pool, id int64, params UpdateTodoParams) error {
	tag, err := db.Exec(ctx, "UPDATE todo SET title = $1, status = $2 WHERE id = $3", params.Title, params.Status, id)
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
	ID     int64
	IDs    []int64
	Title  string
	Status string

	// required
	Limit  int
	Offset int
}

func (s *Store) FindTodos(ctx context.Context, params FindTodosParams) ([]Todo, error) {
	todos, err := findTodos(ctx, s.db, params)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func findTodos(ctx context.Context, db *pgxpool.Pool, params FindTodosParams) ([]Todo, error) {
	var todos []Todo

	rows, err := db.Query(ctx, "select * from todo limit $1 offset $2", params.Limit, params.Offset)
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

func (s *Store) DeleteTodo(ctx context.Context, id int64) error {
	err := deleteTodo(ctx, s.db, id)
	if err != nil {
		return err
	}

	return nil
}

func deleteTodo(ctx context.Context, db *pgxpool.Pool, id int64) error {
	tag, err := db.Exec(ctx, "DELETE FROM todo WHERE id = $1", id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no row deleted")
	}

	return nil
}
