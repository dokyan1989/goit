package postgres

import (
	"database/sql"
	_ "embed"
	"time"

	_ "github.com/go-sql-driver/mysql"
	zlog "github.com/rs/zerolog/log"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

func New(connString string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(time.Hour)

	loggerAdapter := zerologadapter.New(zlog.Logger)
	db = sqldblogger.OpenDriver(connString, db.Driver(), loggerAdapter)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
