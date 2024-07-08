package sqlite

import (
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
	zlog "github.com/rs/zerolog/log"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	loggerAdapter := zerologadapter.New(zlog.Logger)
	db = sqldblogger.OpenDriver(dsn, db.Driver(), loggerAdapter)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
