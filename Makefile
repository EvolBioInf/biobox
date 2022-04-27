progs = al blast2dot bwt cres cutSeq dnaDist drag drawf drawGenes drawKt drawSt fasta2tab geco genTree getSeq huff hut histogram kerror keyMat midRoot maf mtf \
mum2plot mutator naiveMatcher nj num2char numAl pam plotLine plotSeg plotTree pps randomizeSeq ranDot ranseq rep2plot \
repeater revComp rpois sass sblast sequencer shustring simNorm simOrf testMeans translate travTree upgma var watterson \
wrapSeq
packs = util

all:
	test -d bin || mkdir bin
	for pack in $(packs); do \
		make -C $$pack; \
	done
	printf "" > progs.txt
	for prog in $(progs); do \
		make -C $$prog; \
		echo $$prog >> progs.txt; \
		cp $$prog/$$prog bin; \
	done
.PHONY: doc
doc:
	make -C doc
clean:
	for prog in $(progs) $(packs) doc; do \
		make clean -C $$prog; \
	done
	rm -f bin/* progs.txt
test:
	for prog in $(progs) $(packs); do \
		make test -C $$prog; \
	done
