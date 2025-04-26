package cmd

import (
	"goforge/pkg/docs"

	"github.com/urfave/cli/v2"
)

// DocsCommand returns the CLI command for documentation generation.
func DocsCommand() *cli.Command {
	return &cli.Command{
		Name:    "docs",
		Aliases: []string{"doc"},
		Usage:   "Generate documentation",
		Subcommands: []*cli.Command{
			{
				Name:  "api",
				Usage: "Generate API documentation",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "api-docs",
						Usage:   "Output directory for API documentation",
					},
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"f"},
						Value:   "html",
						Usage:   "Output format (html, markdown)",
					},
				},
				Action: func(c *cli.Context) error {
					path := c.Args().First()
					if path == "" {
						path = "."
					}
					return docs.GenerateAPIDoc(path, c.String("output"), c.String("format"))
				},
			},
			{
				Name:  "user",
				Usage: "Generate user documentation",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "user-docs",
						Usage:   "Output directory for user documentation",
					},
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"f"},
						Value:   "html",
						Usage:   "Output format (html, markdown)",
					},
				},
				Action: func(c *cli.Context) error {
					path := c.Args().First()
					if path == "" {
						path = "."
					}
					return docs.GenerateUserDoc(path, c.String("output"), c.String("format"))
				},
			},
		},
	}
}
