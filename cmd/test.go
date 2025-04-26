package cmd

import (
	"goforge/pkg/testing"

	"github.com/urfave/cli/v2"
)

// TestCommand returns the CLI command for test management.
func TestCommand() *cli.Command {
	return &cli.Command{
		Name:  "test",
		Usage: "Test management utilities",
		Subcommands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate tests for Go files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "Output directory for generated tests (defaults to same directory as source)",
					},
					&cli.BoolFlag{
						Name:    "table",
						Aliases: []string{"t"},
						Usage:   "Generate table-driven tests",
					},
				},
				Action: func(c *cli.Context) error {
					path := c.Args().First()
					if path == "" {
						return cli.Exit("Please specify a file or directory to generate tests for", 1)
					}
					output := c.String("output")
					table := c.Bool("table")
					return testing.GenerateTests(path, output, table)
				},
			},
			{
				Name:  "coverage",
				Usage: "Analyze test coverage",
				Flags: []cli.Flag{
					&cli.Float64Flag{
						Name:    "threshold",
						Aliases: []string{"t"},
						Value:   80.0,
						Usage:   "Coverage threshold percentage",
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "coverage.html",
						Usage:   "Output file for coverage report",
					},
				},
				Action: func(c *cli.Context) error {
					path := c.Args().First()
					if path == "" {
						path = "."
					}
					return testing.AnalyzeCoverage(path, c.Float64("threshold"), c.String("output"))
				},
			},
		},
	}
}
