package helloweb

import (
	"embed"
	"io/fs"
)

//go:embed all:tmpl all:static
var embedFS embed.FS

func mustSubFS(sub string) fs.FS {
	f, err := fs.Sub(embedFS, sub)
	if err != nil {
		panic(err)
	}
	return f
}
