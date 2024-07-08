package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	pgxzlog "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	zlog "github.com/rs/zerolog/log"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

func NewPGX(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxzlog.NewLogger(zlog.Logger),
		LogLevel: tracelog.LogLevelInfo,
	}

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewStdDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("open db (connection string = %s): %w", connString, err)
	}

	loggerAdapter := zerologadapter.New(zlog.Logger)
	db = sqldblogger.OpenDriver(connString, db.Driver(), loggerAdapter)

	return db, nil
}
