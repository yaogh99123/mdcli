package main

import "embed"

//go:embed md_source/dist/data.json
var dataJSON []byte

//go:embed md_source/command/*.md
var commandFS embed.FS
