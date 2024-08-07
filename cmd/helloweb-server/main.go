package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/dokyan1989/goit/app/helloweb"
	"github.com/dokyan1989/goit/internal/todo"
	"github.com/dokyan1989/goit/misc/database/postgres"
	"github.com/dokyan1989/goit/misc/envar"
)

// initializes configuration options using the following precedence order
// (below item take precedence over the item above it)
// 1. default
// 2. environment variables
// 3. flags
var (
	env   = envar.GetString("HELLO_ENV", "local")
	host  = envar.GetString("HELLO_HOST", "localhost")
	port  = envar.GetInt("HELLO_PORT", 3000)
	dburl = envar.GetString("HELLO_DB_URL", "")
)

func init() {
	flag.StringVar(&env, "env", env, "running server environment")
	flag.StringVar(&host, "host", host, "server host")
	flag.IntVar(&port, "port", port, "server port")
	flag.StringVar(&dburl, "dburl", dburl, "database connection string")
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()

	// log config values on LOCAL environment
	if strings.ToLower(env) == "local" {
		printConfig()
	}

	if dburl == "" {
		return errors.New("database connection string must be provided")
	}

	// init db
	db, err := postgres.NewPGX(context.Background(), dburl)
	if err != nil {
		return err
	}

	opts := []helloweb.Option{
		helloweb.WithEnv(env),
		helloweb.WithHost(host),
		helloweb.WithPort(port),
	}
	// init server with all dependency
	svr, err := helloweb.NewServer(
		todo.NewStore(db),
		opts...,
	)
	if err != nil {
		return err
	}

	// cleanup all resources before server shutdown
	cleanup := func() {
		db.Close()
	}

	// start server
	svr.Serve(cleanup)
	return nil
}

func printConfig() {
	fmt.Printf("[App configuration] env:%v, host:%v, port:%v, dburl:%v\n",
		env,
		host,
		port,
		dburl,
	)
}
