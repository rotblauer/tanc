### To reproduce the 1000 genomes plots


**Download autosomal vcfs**

```bash

for i in {1..22}
do
   wget ftp://ftp-trace.ncbi.nih.gov/1000genomes/ftp/release/20130502/ALL.chr$i.phase3_shapeit2_mvncall_integrated_v5a.20130502.genotypes.vcf.gz
done

```

**Concatenate vcfs to a single file**

```bash
bcftools concat -n -Oz -o all.vcf.gz *.genotypes.vcf.gz
```

**Create list of ancestry informative markers (AIMS)***

1. Download `Elhaik2013_-_Table_S2.xlsx` from the supplemental materials of 
   [*The GenoChip: A New Tool for Genetic Anthropology*](https://www.ncbi.nlm.nih.gov/pmc/articles/PMC3673633/)
2. Extract rs IDs to `aims.txt`
    
**Run tance**

```bash
tance vc --vcf all.vcf.gz -ids aims.txt -o ~/tanceTest/

```