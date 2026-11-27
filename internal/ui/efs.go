package ui

import (
	"embed"
	"io/fs"
)

//go:embed templates static
var files embed.FS

var (
	HTMLFiles   = sub(files, "templates")
	StaticFiles = sub(files, "static")
)

func sub(f embed.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
