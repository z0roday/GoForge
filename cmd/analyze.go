package cmd

import (
	"goforge/pkg/analyzer"

	"github.com/urfave/cli/v2"
)

// AnalyzeCommand returns the CLI command for analyzing Go code.
func AnalyzeCommand() *cli.Command {
	return &cli.Command{
		Name:    "analyze",
		Aliases: []string{"a"},
		Usage:   "Analyze your Go project structure and code quality",
		Subcommands: []*cli.Command{
			{
				Name:  "structure",
				Usage: "Analyze project structure and architecture",
				Action: func(c *cli.Context) error {
					path := c.Args().First()
					if path == "" {
						path = "."
					}
					return analyzer.AnalyzeStructure(path)
				},
			},
			{
				Name:  "quality",
				Usage: "Analyze code quality and suggest improvements",
				Action: func(c *cli.Context) error {
					path := c.Args().First()
					if path == "" {
						path = "."
					}
					return analyzer.AnalyzeQuality(path)
				},
			},
		},
	}
}
