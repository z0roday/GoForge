package cmd

import (
	"goforge/pkg/container"

	"github.com/urfave/cli/v2"
)

// ContainerCommand returns the CLI command for container operations.
func ContainerCommand() *cli.Command {
	return &cli.Command{
		Name:    "container",
		Aliases: []string{"cont"},
		Usage:   "Generate and manage container configurations",
		Subcommands: []*cli.Command{
			{
				Name:  "dockerfile",
				Usage: "Generate a Dockerfile for your Go application",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "Dockerfile",
						Usage:   "Output file path",
					},
					&cli.StringFlag{
						Name:    "base",
						Aliases: []string{"b"},
						Value:   "golang:alpine",
						Usage:   "Base Docker image",
					},
				},
				Action: func(c *cli.Context) error {
					path := c.Args().First()
					if path == "" {
						path = "."
					}
					return container.GenerateDockerfile(path, c.String("output"), c.String("base"))
				},
			},
			{
				Name:  "kubernetes",
				Usage: "Generate Kubernetes manifests",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "kubernetes",
						Usage:   "Output directory for Kubernetes manifests",
					},
					&cli.StringFlag{
						Name:    "image",
						Aliases: []string{"i"},
						Usage:   "Docker image to use in Kubernetes manifests",
					},
				},
				Action: func(c *cli.Context) error {
					path := c.Args().First()
					if path == "" {
						path = "."
					}
					return container.GenerateKubernetesManifests(path, c.String("output"), c.String("image"))
				},
			},
		},
	}
}
