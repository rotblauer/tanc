package main

import (
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "tance"
	app.Usage = "Generate T-SNE plot from genotypes to visualize ancestry"
	app.Version = "v0.0.1"
	var threads int

	app.Flags = []cli.Flag {
		cli.IntFlag{
			Name:        "threads, t",
			Usage:       "number of threads to use `INT` ",
			Value:       4,
			Destination: &threads,
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:        "tsne",
			Aliases:     []string{"t"},
			Category:    "tsne",
			Usage:       "generate tsne from genotype data",
			UsageText:   "tsne - generates tsne from genotype data",
			Description: "sss",
			ArgsUsage:   "[arrgh]",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "forever, forevvarr"},
			},
			Subcommands: cli.Commands{
				cli.Command{
					Name:   "wop",
					//Action: wopAction,
				},
			},
		},
	}

	app.Run(os.Args)
}
