package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "UniqueFile"
	app.Version = Version
	app.Usage = ""
	app.Author = "Masa Jobara"
	app.Email = "wolf.masa@gmail.com"
	app.Commands = Commands

	app.Run(os.Args)
}
