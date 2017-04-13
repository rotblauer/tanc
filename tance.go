package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

//TODO
// read vcf to float[][]
// tsneGO
// dump results
// plot?

func main() {
	app := cli.NewApp()
	app.Name = "tance"
	app.Usage = "Generates t-Distributed Stochastic Neighbor Embedding (t-SNE) from genotype data"
	app.Version = "v0.0.1"
	var threads int
	var vcf string

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "threads, t",
			Usage:       "number of threads to use `INT` ",
			Value:       4,
			Destination: &threads,
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:    "vcf-compute",
			Aliases: []string{"vc"},
			//Category: "compute",
			Usage:       "compute t-SNE from vcf input ",
			//Description: "compute (t-SNE) from genotypes in vcf file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "vcf, f",
					Usage:       "vcf `FILE` to use ",
					Destination: &vcf,
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("using vcf", vcf)
				fmt.Println("using threads", threads)

				return nil
			},
		},
	}
	app.Run(os.Args)

}
