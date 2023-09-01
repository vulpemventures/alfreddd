package main

import (
	"os"

	"github.com/urfave/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	folderFlag  = "folder"
	projectFlag = "gh-project"
	moduleFlag  = "go-module"
	appFlag     = "app-name"
	apiFlag     = "no-api-spec"
)

func main() {
	app := cli.NewApp()

	app.Version = formatVersion()
	app.Name = "alfreddd"
	app.Usage = "qualified assistant for kick-starting DDD projects"
	app.Action = makeBoilerplate
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  folderFlag,
			Usage: "the folder where I should create the local repository",
		},
		cli.StringFlag{
			Name:  projectFlag,
			Usage: "the name of the GH project in the form <org>/<repo>",
		},
		cli.StringFlag{
			Name:  moduleFlag,
			Usage: "the name of the GO module",
		},
		cli.StringFlag{
			Name:  appFlag,
			Usage: "the name of the application",
		},
		cli.BoolFlag{
			Name:  apiFlag,
			Usage: "skip creating the api-spec folder and buf initialization",
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
