package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	zlog "github.com/rs/zerolog/log"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Container struct {
	container  *mysql.MySQLContainer
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

	mysqlC, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:latest"),
		mysql.WithDatabase(c.dbname),
		mysql.WithUsername(c.dbuser),
		mysql.WithPassword(c.dbpassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog(`port: 3306  MySQL Community Server - GPL`).WithStartupTimeout(5*time.Minute),
		),
		testcontainers.WithEnv(map[string]string{
			"TZ": "UTC",
		}),
		// testcontainers.WithStartupCommand(testcontainers.NewRawCommand([]string{
		// 	"--innodb_doublewrite=OFF",
		// 	"--innodb_file_per_table=OFF",
		// })),
		// https://stackoverflow.com/questions/23631394/mysql-configuration-for-running-tests-faster
		// https://dev.mysql.com/doc/refman/8.0/en/innodb-parameters.html
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Cmd: []string{
					"--innodb_doublewrite=OFF",
					"--innodb_file_per_table=OFF",
				},
			},
		}),
	)
	if err != nil {
		log.Printf("failed to start container: %s\n", err)
		return nil, err
	}

	containerHost, err := mysqlC.Host(ctx)
	if err != nil {
		log.Printf("failed to get container host: %s\n", err)
		return nil, err
	}

	mappedPort, err := mysqlC.MappedPort(ctx, "3306/tcp")
	if err != nil {
		log.Printf("failed to get container port: %s\n", err)
		return nil, err
	}

	c.container = mysqlC
	// https://go.dev/doc/tutorial/database-access#get_handle
	// %s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC&multiStatements=true
	c.connString = (&mysqldriver.Config{
		User:            c.dbuser,
		Passwd:          c.dbpassword,
		Net:             "tcp",
		Addr:            fmt.Sprintf("%s:%s", containerHost, mappedPort.Port()),
		DBName:          c.dbname,
		ParseTime:       true,
		Loc:             time.UTC,
		MultiStatements: true,
	}).FormatDSN()

	return c, nil
}

func (c *Container) GetConnString() string {
	return c.connString
}

func (c *Container) Terminate(ctx context.Context) error {
	return c.container.Terminate(ctx)
}

func (c *Container) OpenDB(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("mysql", c.connString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	loggerAdapter := zerologadapter.New(zlog.Logger)
	db = sqldblogger.OpenDriver(c.connString, db.Driver(), loggerAdapter)

	// force a connection and test that it worked
	i := 0
	for {
		i++
		err := db.Ping()
		if err == nil {
			break
		}

		// TestContainerDBMaxConnectRetry
		if i > 1000 {
			return nil, fmt.Errorf("more than 1000 times fail connect to DB, error: %v", err)
		}
	}

	return db, nil
}

func (c *Container) Migrate(ctx context.Context, sourceURL string) error {
	m, err := migrate.New(sourceURL, "mysql://"+c.connString)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (c *Container) DropAllTables(ctx context.Context) error {
	db, err := c.OpenDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	// prepare sql queries for dropping all tables.
	// https://stackoverflow.com/questions/3476765/mysql-drop-all-tables-ignoring-foreign-keys
	s := "SELECT concat('DROP TABLE IF EXISTS `', table_name, '`;')"
	s += " FROM information_schema.tables"
	s += " WHERE table_schema = ?;"
	rows, err := db.Query(s, c.dbname)
	if err != nil {
		return err
	}
	defer rows.Close()

	dropSQLs := []string{"SET FOREIGN_KEY_CHECKS = 0;"}
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return err
		}

		dropSQLs = append(dropSQLs, s)
	}
	dropSQLs = append(dropSQLs, "SET FOREIGN_KEY_CHECKS = 1;")

	if err := rows.Err(); err != nil {
		return err
	}

	// execute all drop table sql queries
	_, err = db.Exec(strings.Join(dropSQLs, ""))
	if err != nil {
		return err
	}

	return nil
}

func (c *Container) TruncateAllTables(ctx context.Context) error {
	db, err := c.OpenDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	// prepare sql queries for truncating all tables.
	// https://stackoverflow.com/questions/5452760/how-to-truncate-a-foreign-key-constrained-table
	s := "SELECT concat('TRUNCATE TABLE `', table_name, '`;')"
	s += " FROM information_schema.tables"
	s += " WHERE table_schema = ?;"
	rows, err := db.Query(s, c.dbname)
	if err != nil {
		return err
	}
	defer rows.Close()

	truncateSQLs := []string{"SET FOREIGN_KEY_CHECKS = 0;SET UNIQUE_CHECKS = 0;"}
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return err
		}

		truncateSQLs = append(truncateSQLs, s)
	}
	truncateSQLs = append(truncateSQLs, "SET FOREIGN_KEY_CHECKS = 1;SET UNIQUE_CHECKS = 1;")

	if err := rows.Err(); err != nil {
		return err
	}

	// execute all truncate table sql queries
	_, err = db.Exec(strings.Join(truncateSQLs, ""))
	if err != nil {
		return err
	}

	return nil
}
