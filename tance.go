package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/biogo/hts/bgzf"
	"github.com/brentp/vcfgo"
	"github.com/rotblauer/goTsne/Utils"
	"github.com/rotblauer/tsne4go"
	"github.com/urfave/cli"
	"os"
	"strconv"
)

func extractGenotypes(variant *vcfgo.Variant, header *vcfgo.Header) []float64 {
	var s []float64
	header.ParseSamples(variant)
	t := variant.Samples
	for _, gContext := range t {
		gt := 0
		for _, aCount := range gContext.GT {
			gt += aCount
		}
		if gt < 0 {
			fmt.Println(gContext)
		}
		s = append(s, float64(gt))

	}
	return s

}

func loadData(vcf string, idFile string, outDir string, limit bool, iter int, temp int, perplexity float64, epsilon float64) {
	if limit {

		rsIds := Utils.LoadRsId(idFile)
		fmt.Printf("%d total rsIds loaded\n", len(rsIds))

		run(vcf, rsIds, outDir, iter, temp, perplexity, epsilon)
	} else {
		run(vcf, make(map[string]string), outDir, iter, temp, perplexity, epsilon)

	}
}

func transpose(a [][]float64) [][]float64 {
	n := len(a[0])
	b := make([][]float64, n)
	for i := 0; i < n; i++ {
		b[i] = make([]float64, len(a))
		for j := 0; j < len(a); j++ {
			b[i][j] = a[j][i]
		}
	}
	return b
}

//http://distill.pub/2016/misread-tsne/
func run(vcf string, rsIds map[string]string, outDir string, iter int, temp int, perplexity float64, epsilon float64) {
	f, _ := os.Open(vcf)
	//TODO non-gzip based on ext
	//TODO root based output
	r, err := gzip.NewReader(f)
	rdr, err := vcfgo.NewReader(r, true)
	if err != nil {
		panic(err)
	}
	outVCF := outDir + "tance_variants_used.vcf.gz"
	fmt.Println("writing variants used to " + outVCF)

	os.MkdirAll(outDir, os.ModePerm)
	o, err := os.Create(outVCF)

	bgzfw := bgzf.NewWriter(bufio.NewWriter(o), 4)
	wtr, err := vcfgo.NewWriter(bgzfw, rdr.Header)
	defer o.Close()
	defer bgzfw.Flush()

	var genotypeMatrix [][]float64
	num := 0
	numUsed := 0
	for {
		variant := rdr.Read()
		if variant == nil {
			break
		}
		if num%10000 == 0 {
			fmt.Printf("%d total variants scanned\n", num)
			fmt.Printf("%d total variants kept\n", numUsed)

		}
		if _, ok := rsIds[variant.Id()]; ok || len(rsIds) == 0 {
			wtr.WriteVariant(variant)
			numUsed++
			genotypeMatrix = append(genotypeMatrix, extractGenotypes(variant, rdr.Header))
		}
		num++

	}
	samples := make([]interface{}, len(rdr.Header.SampleNames))
	for i, v := range rdr.Header.SampleNames {
		samples[i] = v
	}

	tsne := tsne4go.New(Utils.GenotypeDistancer{Matrix: transpose(genotypeMatrix)}, samples, perplexity, tsne4go.ToleranceDefault)
	for i := 0; i < iter; i++ {
		tsne.Step(epsilon)
		fmt.Println(i)
		if i%temp == 0 && temp > 0 {
			dumpCurrent(outDir+"tance_tsne_rep_"+strconv.Itoa(i)+".txt", tsne, rdr)

		}
	}

	dumpCurrent(outDir+"tance_tsne_final.txt", tsne, rdr)
}
func dumpCurrent(out string, tsne *tsne4go.TSne, rdr *vcfgo.Reader) {
	tsneOut, err := os.Create(out)
	wTsne := bufio.NewWriter(tsneOut)
	defer tsneOut.Close()
	for sampIndex, point := range tsne.Solution {
		_, err = fmt.Fprintf(wTsne, "%s", rdr.Header.SampleNames[sampIndex])
		check(err)
		for _, coord := range point {
			fmt.Fprint(wTsne, "\t")
			fmt.Fprint(wTsne, coord)
		}
		fmt.Fprint(wTsne, "\n")
	}
	wTsne.Flush()
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "tance"
	app.Usage = "Generates t-Distributed Stochastic Neighbor Embedding (t-SNE) from genotype data"
	app.Version = "v0.0.1"
	var vcf string
	var idFile string
	var outDir string
	var iter int
	var temp int
	var perplexity float64
	var epsilon float64
	var limit = false

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "threads, t",
			Usage:  "number of threads to use `INT` ",
			Value:  4,
			EnvVar: "GOMAXPROCS",
		},
	}
	app.Commands = []cli.Command{
		{
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
					Name:        "ids",
					Usage:       "idFile `FILE` only variant IDs within this file will be used. If not supplied, no input filtering will occur ",
					Destination: &idFile,
				},
				cli.StringFlag{
					Name:        "outDir, o",
					Usage:       "outDir `DIR` output directory ",
					Destination: &outDir,
				},
				cli.IntFlag{
					Name:        "iterations, i",
					Usage:       "number of iterations for t-SNE `INT` ",
					Value:       10000,
					Destination: &iter,
				},
				cli.IntFlag{
					Name:        "report, r",
					Usage:       "number of iterations for temporary solutions `INT` ",
					Value:       1000,
					Destination: &temp,
				},
				cli.Float64Flag{
					Name:        "perplexity, p",
					Usage:       "The performance of SNE is fairly robust to changes in the perplexity, and typical values are between 5 and 50. `FLOAT` ",
					Value:       tsne4go.PerplexityDefault,
					Destination: &perplexity,
				},
				cli.Float64Flag{
					Name:        "epsilon, e",
					Usage:       "A learning rate (often called “epsilon”)  `FLOAT` ",
					Value:       tsne4go.EpsilonDefault,
					Destination: &epsilon,
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("using vcf", vcf)
				if idFile != "" {
					fmt.Println("using limiting ID file", idFile)
					limit = true

				}

				loadData(vcf, idFile, outDir, limit, iter, temp, perplexity, epsilon)
				return nil
			},
		},
	}
	app.Run(os.Args)

}
