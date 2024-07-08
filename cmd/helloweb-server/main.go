package main

import (
	"flag"
	"fmt"
	"strings"

	web "github.com/dokyan1989/goit/app/helloweb"
	"github.com/dokyan1989/goit/misc/envar"
)

// initializes configuration options using the following precedence order
// (below item take precedence over the item above it)
// 1. default
// 2. environment variables
// 3. flags
var (
	env  = envar.GetString("HELLO_ENV", "local")
	host = envar.GetString("HELLO_HOST", "localhost")
	port = envar.GetInt("HELLO_PORT", 3000)
)

func init() {
	flag.StringVar(&env, "env", env, "running server environment")
	flag.StringVar(&host, "host", host, "server host")
	flag.IntVar(&port, "port", port, "server port")
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

	opts := []web.Option{
		web.WithEnv(env),
		web.WithHost(host),
		web.WithPort(port),
	}
	// init server with all dependency
	svr, err := web.NewServer(opts...)
	if err != nil {
		return err
	}

	// start server
	svr.Serve()
	return nil
}

func printConfig() {
	fmt.Printf("[App configuration] env:%v, host:%v, port:%v\n", env, host, port)
}
