package main

import (
	"context"

	"github.com/dokyan1989/goit/app/grpc"
	todoGRPCSvc "github.com/dokyan1989/goit/internal/todo/impl/grpc"
	"github.com/dokyan1989/goit/util/database/postgres"
)

func main() {
	connString := "postgres://username:password@localhost:5432/sample?sslmode=disable"
	db, err := postgres.New(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	todoAdapter := todoGRPCSvc.New(db)

	svr := grpc.NewServer(todoAdapter)
	svr.Serve()
}
