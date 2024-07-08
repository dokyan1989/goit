package main

import (
	"context"

	"github.com/dokyan1989/goit/app/api"
	todoAPISvc "github.com/dokyan1989/goit/internal/todo/impl/api"
	"github.com/dokyan1989/goit/util/database/postgres"
)

func main() {
	connString := "postgres://username:password@localhost:5432/sample?sslmode=disable"
	db, err := postgres.New(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	todoSvc := todoAPISvc.New(db)

	svr, err := api.NewServer(todoSvc)
	if err != nil {
		panic(err)
	}

	svr.Serve()
}
