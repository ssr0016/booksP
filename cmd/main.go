package main

import (
	"os"
	"practice/cmd/app"
)

func main() {
	command := app.NewCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
