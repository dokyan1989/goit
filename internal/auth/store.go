package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

type CreateUserParams struct {
	UserName  string
	Password  string
	Email     string
	FirstName string
	LastName  string
	Status    string
}

func (s *Store) CreateUser(ctx context.Context, p CreateUserParams) (int64, error) {
	// hash and salt password
	hashedPwd, err := hashAndSalt([]byte(p.Password))
	if err != nil {
		return 0, fmt.Errorf("hash and salt password: %w", err)
	}

	p.Password = hashedPwd
	id, err := createUser(ctx, s.db, p)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Store) FindUserByID(ctx context.Context, id int64) (User, error) {
	return findUserByID(ctx, s.db, id)
}

func (s *Store) FindUserByNameAndPassword(ctx context.Context, name, password string) (User, error) {
	user, err := findUserByName(ctx, s.db, name)
	if err != nil {
		return User{}, err
	}

	err = comparePasswords(user.Password, password)
	if err != nil {
		return User{}, errors.New("password is incorrect")
	}

	return user, nil
}

func (s *Store) CreateUserRole(ctx context.Context, userID, roleID int64) (int64, error) {
	sql := `INSERT INTO user_role (user_id, role_id) VALUES ($1, $2) RETURNING user_role_id`

	var id int64
	row := s.db.QueryRow(ctx, sql, userID, roleID)
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("no row inserted")
		}
		return 0, fmt.Errorf("row.Scan() error: %w", err)
	}

	return id, nil
}

type CreateRoleParams struct {
	RoleName    string
	Description string
}

func (s *Store) CreateRole(ctx context.Context, p CreateRoleParams) (int64, error) {
	sql := `INSERT INTO role (role_name, description) VALUES ($1, $2) RETURNING role_id`

	var id int64
	row := s.db.QueryRow(ctx, sql, p.RoleName, p.Description)
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("no row inserted")
		}
		return 0, fmt.Errorf("row.Scan() error: %w", err)
	}

	return id, nil
}

func (s *Store) FetchUserWithRole(ctx context.Context, userID int64) (UserWithRoles, error) {
	user, err := s.FindUserByID(ctx, userID)
	if err != nil {
		return UserWithRoles{}, err
	}

	var roles []Role
	sql := `
SELECT r.role_id, r.role_name, r.description, r.created_at, r.updated_at, r.status
FROM user_role ur JOIN role r ON ur.role_id = r.role_id 
WHERE ur.user_id = $1`

	err = pgxscan.Select(ctx, s.db, &roles, sql, user.UserID)
	if err != nil {
		return UserWithRoles{}, err
	}

	return UserWithRoles{User: user, Roles: roles}, nil
}

// https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72
func hashAndSalt(pwd []byte) (string, error) {
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func createUser(ctx context.Context, db *pgxpool.Pool, p CreateUserParams) (int64, error) {
	var sql = `
INSERT INTO "user"(user_name, password, email, first_name, last_name, status) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id`

	var id int64
	row := db.QueryRow(ctx, sql,
		p.UserName,
		p.Password,
		p.Email,
		p.FirstName,
		p.LastName,
		p.Status,
	)

	if err := row.Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("no row inserted")
		}
		return 0, fmt.Errorf("row.Scan() error: %w", err)
	}

	return id, nil
}

func findUserByID(ctx context.Context, db *pgxpool.Pool, id int64) (User, error) {
	var user User
	var sql = `
SELECT user_id, user_name, password, email, first_name, last_name, created_at, updated_at, last_login_at, status 
FROM "user"
WHERE user_id = $1`

	if err := pgxscan.Get(ctx, db, &user, sql, id); err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Store) UpdateUserStatus(ctx context.Context, id int64, status string) error {
	err := updateUserStatus(ctx, s.db, id, status)
	if err != nil {
		return err
	}

	return nil
}

func updateUserStatus(ctx context.Context, db *pgxpool.Pool, id int64, status string) error {
	tag, err := db.Exec(ctx, `UPDATE "user" SET status = $1 WHERE user_id = $2`, status, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("no row updated")
	}

	return nil
}
func findUserByName(ctx context.Context, db *pgxpool.Pool, name string) (User, error) {
	var user User
	var sql = `
SELECT user_id, user_name, password, email, first_name, last_name, created_at, updated_at, last_login_at, status 
FROM "user"
WHERE user_name = $1`

	if err := pgxscan.Get(ctx, db, &user, sql, name); err != nil {
		return User{}, err
	}

	return user, nil
}

func comparePasswords(hashedPwd, plainPwd string) error {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	bHash, bPlain := []byte(hashedPwd), []byte(plainPwd)

	err := bcrypt.CompareHashAndPassword(bHash, bPlain)
	if err != nil {
		return err
	}

	return nil
}
