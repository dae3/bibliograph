//go:build noembedspa

package main

import (
	"io/fs"
	"os"
)

var SpaFS fs.FS

func init() {
	SpaFS = os.DirFS("client")
}
