package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "dockerboot"
	app.Version = Version
	app.Usage = ""
	app.Author = "takoyaki"
	app.Email = "info@thun.ws"
	app.Commands = Commands

	app.Run(os.Args)
}
