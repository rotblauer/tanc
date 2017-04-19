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
   1. Available (as of 4-18-17) [here](https://oup.silverchair-cdn.com/oup/backfile/Content_public/Journal/gbe/5/5/10.1093_gbe_evt066/1/evt066_Supplementary_Data.zip?Expires=1492662814&Signature=HSnHR-lWbQnt302x2h-9Srd-ECUHY9gRizGaeNE2N2UfhY47hV8M9rUcCkKxXYH-SmszdZnAzZBZ~1LkB2zKvLRQrnnv9wY5Q76UtEhOkJAYkslqoArqt53-YcHhLZgcmY1JPUBudtS~XIAKlMVcCWtMqJo6d8IEFgXDAxrAJxwbhCzhbnc-wdcG2fAq5Fd5GXDZiHRBsvkKGSdTiqL-xL5w4L2G6iXqzKi7hUc2dUIpUpTnjEj-BxNRFLX-W40wDIazZnuoRi2kNGau-s22oFYXLnI317iygJd7IIymFFCeBRL7Ep6pVNYjSIdfaDhCT1Wy-ObgEwo8zLNO1OAdpA__&Key-Pair-Id=APKAIUCZBIA4LVPAVW3Q)
2. Extract IDs to `aims.txt`
    
**Run tance**

```bash
tance vc --vcf all.vcf.gz -ids aims.txt -o ~/tanceTest/

```