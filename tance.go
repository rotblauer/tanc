package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/brentp/vcfgo"
	"github.com/rotblauer/goTsne/Utils"
	"github.com/urfave/cli"
	"os"
)

func loadData(vcf string, idFile string, outDir string) {
	rsIds := Utils.LoadRsId(idFile)
	fmt.Printf("%d total rsIds loaded\n", len(rsIds))

	readVCF(vcf, rsIds, outDir)
}

func readVCF(vcf string, rsIds map[string]string, outDir string) {
	f, _ := os.Open(vcf)
	r, err := gzip.NewReader(f)
	rdr, err := vcfgo.NewReader(r, false)
	if err != nil {
		panic(err)
	}
	fmt.Print("writing to " + outDir + "tance_variants_used.vcf")
	//o, _ := os.Open()
	o, err := os.Create(outDir + "tance_variants_used.vcf")
	w := bufio.NewWriter(o)
	wtr, err := vcfgo.NewWriter(w, rdr.Header)
	defer o.Close()

	num := 0
	numUsed := 0
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
		if numUsed+1%100 == 0 {
			break
		}
		if _, ok := rsIds[variant.Id_]; ok {
			wtr.WriteVariant(variant)
			numUsed++
		}
	}
	w.Flush()

}

//TODO
// read vcf
// add genotypes to float to float[][]
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
	var idFile string
	var outDir string

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
					Name:        "ids, i",
					Usage:       "idFile `FILE` only variant IDs within this file will be used ",
					Destination: &idFile,
				},
				cli.StringFlag{
					Name:        "outDir, o",
					Usage:       "outDir `DIR` output directory ",
					Destination: &outDir,
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("using vcf", vcf)
				fmt.Println("using threads", threads)
				loadData(vcf, idFile, outDir)
				return nil
			},
		},
	}
	app.Run(os.Args)

}
