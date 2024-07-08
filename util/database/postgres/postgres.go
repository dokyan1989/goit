package postgres

import (
	"context"
	_ "embed"

	pgxzlog "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	zlog "github.com/rs/zerolog/log"
)

func New(ctx context.Context, connString string) (*pgxpool.Pool, error) {
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
