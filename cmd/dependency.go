package cmd

import (
	"goforge/pkg/dependency"

	"github.com/urfave/cli/v2"
)

// DependencyCommand returns the CLI command for managing dependencies.
func DependencyCommand() *cli.Command {
	return &cli.Command{
		Name:    "dependency",
		Aliases: []string{"dep"},
		Usage:   "Manage project dependencies",
		Subcommands: []*cli.Command{
			{
				Name:  "check",
				Usage: "Check for outdated dependencies",
				Action: func(c *cli.Context) error {
					path := c.Args().First()
					if path == "" {
						path = "."
					}
					return dependency.CheckOutdated(path)
				},
			},
			{
				Name:  "update",
				Usage: "Update dependencies to latest versions",
				Action: func(c *cli.Context) error {
					path := c.Args().First()
					if path == "" {
						path = "."
					}
					return dependency.Update(path)
				},
			},
			{
				Name:  "security",
				Usage: "Check dependencies for security vulnerabilities",
				Action: func(c *cli.Context) error {
					path := c.Args().First()
					if path == "" {
						path = "."
					}
					return dependency.CheckSecurity(path)
				},
			},
		},
	}
}
