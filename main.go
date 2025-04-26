package main

import (
	"log"
	"os"

	"goforge/cmd"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "goforge",
		Usage: "A comprehensive development companion for Go projects",
		Commands: []*cli.Command{
			cmd.AnalyzeCommand(),
			cmd.DependencyCommand(),
			cmd.ProfileCommand(),
			cmd.ContainerCommand(),
			cmd.TestCommand(),
			cmd.DocsCommand(),
			cmd.APICommand(),
			cmd.WebCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
