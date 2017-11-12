#!/usr/bin/env bash 

go run tance.go vc --vcf /Users/Kitty/tmp/tance_variants_used.vcf.gz -ids /Users/Kitty/Go/src/github.com/rotblauer/goTsne/ext/aims.txt -o /Users/Kitty/tmp/tance/



for p in $(seq 0 5 75)
do
for e in $(seq 0 5 75)
do
	go run tance.go -t 4 vc --vcf /Users/Kitty/tmp/tance_variants_used.vcf.gz -o "/Users/Kitty/tmp/tance/iters/perp_"$p"_epsi_"$e -p $p -e $e -i 1000 -r 100
done
done
