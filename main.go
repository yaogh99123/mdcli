package main

import (
	"mdcli/pkg/app"
	"mdcli/pkg/commands"
)

func main() {
	commands.SetEmbeddedData(dataJSON, commandFS)
	app.Run()
}
