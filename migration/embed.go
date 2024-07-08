package migration

import "embed"

var (
	//go:embed todo
	TodoFS embed.FS
	//go:embed auth
	AuthFS embed.FS
)

const (
	PathTodo = "todo"
	PathAuth = "auth"
)
