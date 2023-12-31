//go:build !noembedspa

package main

import (
	"embed"
	"io/fs"
)

//go:embed client-dist/*
var spaFS embed.FS
var SpaFS fs.FS

func init() {
	var err error
	SpaFS, err = fs.Sub(spaFS, "client-dist")
	if err != nil {
		panic("Can't prefix mount embedded SPA filesystem")
	}
}
