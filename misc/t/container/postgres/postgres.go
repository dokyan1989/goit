package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	pgxzlog "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	zlog "github.com/rs/zerolog/log"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
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
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("run container: %w", err)
	}

	containerHost, err := postgresC.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("get container host: %w", err)
	}

	mappedPort, err := postgresC.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, fmt.Errorf("get container port: %w", err)
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

func (c *Container) OpenPGXDB(ctx context.Context) (*pgxpool.Pool, error) {
	dbpool, err := openPGXDB(ctx, c.connString)
	if err != nil {
		return nil, err
	}

	if err := dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping to db: %w", err)
	}

	return dbpool, nil
}

// func (c *Container) Migrate(ctx context.Context, sourceURL string) error {
// 	m, err := migrate.New(sourceURL, c.connString)
// 	if err != nil {
// 		return err
// 	}

// 	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
// 		return err
// 	}

// 	return nil
// }

// Migrate migrates with directory. E.g: file://path/to/root/internal/migration/todo
func (c *Container) Migrate(ctx context.Context, sourceURL string) error {
	db, err := openStdDB(c.connString)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := migratepostgres.WithInstance(db, &migratepostgres.Config{})
	if err != nil {
		return fmt.Errorf("get driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return fmt.Errorf("get instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

// MigrateFS migrates with embedded files
func (c *Container) MigrateFS(ctx context.Context, fs fs.FS, path string) error {
	sourceDriver, err := iofs.New(fs, path)
	if err != nil {
		return fmt.Errorf("iofs: %w", err)
	}

	db, err := openStdDB(c.connString)
	if err != nil {
		return err
	}
	defer db.Close()

	dbDriver, err := migratepostgres.WithInstance(db, &migratepostgres.Config{})
	if err != nil {
		return fmt.Errorf("get driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return fmt.Errorf("get instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

// func (c *Container) Seed(ctx context.Context, sourceURL string) error {
// 	// https://stackoverflow.com/a/78161246
// 	// https://github.com/golang-migrate/migrate/tree/master/database/pgx/v5
// 	m, err := migrate.New(sourceURL, c.connString+"&x-migrations-table=seed_migrations")
// 	if err != nil {
// 		return err
// 	}

// 	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
// 		return err
// 	}

// 	return nil
// }

func (c *Container) Seed(ctx context.Context, sourceURL string) error {
	db, err := openStdDB(c.connString)
	if err != nil {
		return err
	}
	defer db.Close()

	// https://stackoverflow.com/a/78161246
	// https://github.com/golang-migrate/migrate/tree/master/database/pgx/v5
	driver, err := migratepostgres.WithInstance(db, &migratepostgres.Config{
		MigrationsTable: "seed_migrations",
	})
	if err != nil {
		return fmt.Errorf("get driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return fmt.Errorf("get instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

func (c *Container) DropAllTables(ctx context.Context) error {
	dbpool, err := openPGXDB(ctx, c.connString)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer dbpool.Close()

	// prepare sql queries for dropping all tables.
	// https://www.educative.io/answers/how-to-drop-all-the-tables-in-a-postgresql-database
	rows, err := dbpool.Query(ctx, `SELECT 'DROP TABLE IF EXISTS "' || tablename || '" CASCADE;' FROM pg_tables WHERE schemaname = 'public'`)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	dropSqls := []string{}
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return fmt.Errorf("scan row error: %w", err)
		}

		dropSqls = append(dropSqls, s)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %w", err)
	}

	_, err = dbpool.Exec(ctx, strings.Join(dropSqls, ""))
	if err != nil {
		return fmt.Errorf("execute query error: %w", err)
	}

	return nil
}

func (c *Container) TruncateAllTables(ctx context.Context) error {
	dbpool, err := openPGXDB(ctx, c.connString)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer dbpool.Close()

	// prepare sql queries for truncating all tables.
	// https://www.postgresql.org/docs/current/sql-truncate.html
	rows, err := dbpool.Query(ctx, `SELECT 'TRUNCATE TABLE "' || tablename || '" RESTART IDENTITY CASCADE;' FROM pg_tables WHERE schemaname = 'public'`)
	if err != nil {
		return fmt.Errorf("execute query error: %w", err)
	}
	defer rows.Close()

	truncateSqls := []string{}
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return fmt.Errorf("scan row error: %w", err)
		}

		truncateSqls = append(truncateSqls, s)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %w", err)
	}

	_, err = dbpool.Exec(ctx, strings.Join(truncateSqls, ""))
	if err != nil {
		return fmt.Errorf("execute query error: %w", err)
	}

	return nil
}

func openPGXDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse config (connection string = %s): %w", connString, err)
	}

	cfg.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxzlog.NewLogger(zlog.Logger),
		LogLevel: tracelog.LogLevelInfo,
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create new dbpool: %w", err)
	}

	return dbpool, nil
}

func openStdDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("open db (connection string = %s): %w", connString, err)
	}

	loggerAdapter := zerologadapter.New(zlog.Logger)
	db = sqldblogger.OpenDriver(connString, db.Driver(), loggerAdapter)

	return db, nil
}
