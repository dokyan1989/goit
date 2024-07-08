package main

import (
	"context"

	grpc "github.com/dokyan1989/goit/app/hellogrpc"
	"github.com/dokyan1989/goit/internal/todo"
	"github.com/dokyan1989/goit/misc/database/postgres"
)

func main() {
	connString := "postgres://username:password@localhost:5432/sample?sslmode=disable"
	db, err := postgres.NewPGX(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	svr := grpc.NewServer(todo.NewStore(db))
	svr.Serve()
}
