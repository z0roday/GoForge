package cmd

import (
	"goforge/pkg/profiler"

	"github.com/urfave/cli/v2"
)

// ProfileCommand returns the CLI command for profiling Go applications.
func ProfileCommand() *cli.Command {
	return &cli.Command{
		Name:    "profile",
		Aliases: []string{"prof"},
		Usage:   "Profile your Go application",
		Subcommands: []*cli.Command{
			{
				Name:  "cpu",
				Usage: "Profile CPU usage",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "cpu.pprof",
						Usage:   "Output file for CPU profile",
					},
					&cli.IntFlag{
						Name:    "duration",
						Aliases: []string{"d"},
						Value:   30,
						Usage:   "Duration in seconds to run the profile",
					},
				},
				Action: func(c *cli.Context) error {
					target := c.Args().First()
					if target == "" {
						return cli.Exit("Please specify a binary to profile", 1)
					}
					return profiler.CPUProfile(target, c.String("output"), c.Int("duration"))
				},
			},
			{
				Name:  "memory",
				Usage: "Profile memory usage",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "mem.pprof",
						Usage:   "Output file for memory profile",
					},
				},
				Action: func(c *cli.Context) error {
					target := c.Args().First()
					if target == "" {
						return cli.Exit("Please specify a binary to profile", 1)
					}
					return profiler.MemoryProfile(target, c.String("output"))
				},
			},
			{
				Name:  "visualize",
				Usage: "Visualize profile data",
				Action: func(c *cli.Context) error {
					profile := c.Args().First()
					if profile == "" {
						return cli.Exit("Please specify a profile file to visualize", 1)
					}
					return profiler.Visualize(profile)
				},
			},
		},
	}
}
