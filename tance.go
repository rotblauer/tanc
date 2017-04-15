package main

import (
	"compress/gzip"
	"fmt"
	"github.com/brentp/vcfgo"
	"github.com/urfave/cli"
	"os"
	"github.com/rotblauer/goTsne/Utils"
)

func loadData(vcf string, rsIDsFile string) {
	rsIds :=Utils.LoadRsId(rsIDsFile)
	fmt.Printf("%d total rsIds loaded\n", len(rsIds))

	readVCF(vcf, rsIds)
}

func readVCF(vcf string, rsIds map[string]string) {
	f, _ := os.Open(vcf)
	r, err := gzip.NewReader(f)
	rdr, err := vcfgo.NewReader(r, false)
	if err != nil {
		panic(err)
	}
	num := 0
	numUsed :=0
	for {
		variant := rdr.Read()
		if variant == nil {
			break
		}
		num++
		if num%1000 == 0 {
			fmt.Printf("%d total variants scanned\n", num)
			fmt.Printf("%d total variants kept\n", numUsed)

		}
		if _, ok := rsIds[variant.Id_]; ok {
			//do something here
			numUsed++;
		}

		//fmt.Printf("%s\t%d\t%s\t%s\n", variant.Chromosome, variant.Pos, variant.Ref, variant.Alt)

		//fmt.Print(variant.)
		//fmt.Printf("%s", variant.Info("DP").(int) > 10)
		//sample := variant.Samples[0]
		// we can get the PL field as a list (-1 is default in case of missing value)
		//fmt.Println("%s", variant.GetGenotypeField(sample, "PL", -1))
		//_ = sample.DP
	}
}

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
	var rsIDsFile string

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
			Usage: "compute t-SNE from vcf input ",
			//Description: "compute (t-SNE) from genotypes in vcf file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "vcf, f",
					Usage:       "vcf `FILE` to use ",
					Destination: &vcf,
				},
				cli.StringFlag{
					Name:        "rs-ids, rs",
					Usage:       "rs-ids `FILE` only variant IDs within this file will be used ",
					Destination: &rsIDsFile,
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("using vcf", vcf)
				fmt.Println("using threads", threads)
				loadData(vcf, rsIDsFile)
				return nil
			},
		},
	}
	app.Run(os.Args)

}
