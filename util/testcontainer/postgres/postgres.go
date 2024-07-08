package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	pgxzlog "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	zlog "github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Container struct {
	container  *postgres.PostgresContainer
	connString string
	dbname     string
	dbuser     string
	dbpassword string
}

type Option func(*Container)

func WithDBName(dbname string) Option {
	return func(c *Container) {
		c.dbname = dbname
	}
}

func WithDBUser(dbuser string) Option {
	return func(c *Container) {
		c.dbuser = dbuser
	}
}

func WithDBPassword(dbpassword string) Option {
	return func(c *Container) {
		c.dbpassword = dbpassword
	}
}

func defaultContainer() *Container {
	return &Container{
		dbname:     "sample",
		dbuser:     "default",
		dbpassword: "secret",
	}
}

func Run(ctx context.Context, options ...Option) (*Container, error) {
	c := defaultContainer()
	for _, opt := range options {
		opt(c)
	}

	postgresC, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:latest"),
		postgres.WithDatabase(c.dbname),
		postgres.WithUsername(c.dbuser),
		postgres.WithPassword(c.dbpassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		log.Printf("failed to start container: %s\n", err)
		return nil, err
	}

	containerHost, err := postgresC.Host(ctx)
	if err != nil {
		log.Printf("failed to get container host: %s\n", err)
		return nil, err
	}

	mappedPort, err := postgresC.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Printf("failed to get container port: %s\n", err)
		return nil, err
	}

	c.container = postgresC
	c.connString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.dbuser,
		c.dbpassword,
		containerHost,
		mappedPort.Port(),
		c.dbname,
	)
	return c, nil
}

func (c *Container) GetConnString() string {
	return c.connString
}

func (c *Container) Terminate(ctx context.Context) error {
	return c.container.Terminate(ctx)
}

func (c *Container) OpenDB(ctx context.Context) (*pgxpool.Pool, error) {
	pgConfig, err := pgxpool.ParseConfig(c.connString)
	if err != nil {
		zlog.Printf("Unable to parse config: %v\n", err)
		return nil, err
	}

	pgConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxzlog.NewLogger(zlog.Logger),
		LogLevel: tracelog.LogLevelInfo,
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		zlog.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	if err := dbpool.Ping(ctx); err != nil {
		zlog.Printf("Unable to ping to database: %v\n", err)
		return nil, err
	}

	return dbpool, nil
}

func (c *Container) Migrate(ctx context.Context, sourceURL string) error {
	m, err := migrate.New(sourceURL, c.connString)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (c *Container) DropAllTables(ctx context.Context) error {
	// connect to postgres
	dbpool, err := c.OpenDB(ctx)
	if err != nil {
		return err
	}
	defer dbpool.Close()

	// prepare sql queries for dropping all tables.
	// https://www.educative.io/answers/how-to-drop-all-the-tables-in-a-postgresql-database
	rows, err := dbpool.Query(ctx, `SELECT 'DROP TABLE IF EXISTS "' || tablename || '" CASCADE;' FROM pg_tables WHERE schemaname = 'public'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	dropSqls := []string{}
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return err
		}

		dropSqls = append(dropSqls, s)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// execute all drop table sql queries
	_, err = dbpool.Exec(ctx, strings.Join(dropSqls, ""))
	if err != nil {
		return err
	}

	return nil
}

func (c *Container) TruncateAllTables(ctx context.Context) error {
	// connect to postgres
	dbpool, err := c.OpenDB(ctx)
	if err != nil {
		return err
	}
	defer dbpool.Close()

	// prepare sql queries for truncating all tables.
	// https://www.postgresql.org/docs/current/sql-truncate.html
	rows, err := dbpool.Query(ctx, `SELECT 'TRUNCATE TABLE "' || tablename || '" RESTART IDENTITY CASCADE;' FROM pg_tables WHERE schemaname = 'public'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	truncateSqls := []string{}
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return err
		}

		truncateSqls = append(truncateSqls, s)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// execute all truncate table sql queries
	_, err = dbpool.Exec(ctx, strings.Join(truncateSqls, ""))
	if err != nil {
		return err
	}

	return nil
}
