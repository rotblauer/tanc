package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
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

func loadData(vcf string, idFile string, outDir string, iter int, temp int, perplexity float64, epsilon float64) {
	rsIds := Utils.LoadRsId(idFile)
	fmt.Printf("%d total rsIds loaded\n", len(rsIds))

	run(vcf, rsIds, outDir, iter, temp, perplexity, epsilon)
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
	r, err := gzip.NewReader(f)
	rdr, err := vcfgo.NewReader(r, true)
	if err != nil {
		panic(err)
	}
	fmt.Print("writing to " + outDir + "tance_variants_used.vcf")
	//o, _ := os.Open()
	o, err := os.Create(outDir + "tance_variants_used.vcf")
	w := bufio.NewWriter(o)
	wtr, err := vcfgo.NewWriter(w, rdr.Header)
	defer o.Close()

	var genotypeMatrix [][]float64
	num := 0
	numUsed := 0
	for {
		variant := rdr.Read()
		if variant == nil {
			break
		}
		num++
		if num%10000 == 0 {
			fmt.Printf("%d total variants scanned\n", num)
			fmt.Printf("%d total variants kept\n", numUsed)

		}
		if _, ok := rsIds[variant.Id_]; ok {
			wtr.WriteVariant(variant)
			numUsed++
			genotypeMatrix = append(genotypeMatrix, extractGenotypes(variant, rdr.Header))
		}
	}
	w.Flush()
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
	var threads int
	var vcf string
	var idFile string
	var outDir string
	var iter int
	var temp int
	var perplexity float64
	var epsilon float64

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "threads, t",
			Usage:       "number of threads to use `INT` ",
			Value:       4,
			Destination: &threads,
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
					Usage:       "idFile `FILE` only variant IDs within this file will be used ",
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
					Value:       20000,
					Destination: &iter,
				},
				cli.IntFlag{
					Name:        "report, r",
					Usage:       "number of iterations for temporary solutions `INT` ",
					Value:       200,
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
				fmt.Println("using threads", threads)
				loadData(vcf, idFile, outDir, iter, temp, perplexity, epsilon)
				return nil
			},
		},
	}
	app.Run(os.Args)

}
