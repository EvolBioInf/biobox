packs = util
progs = al blast2dot bwt clac cres cutSeq dnaDist drag drawf drawGenes drawKt \
drawSt fasta2tab geco genTree getSeq huff hut histogram kerror keyMat midRoot maf mtf \
mum2plot mutator naiveMatcher nj num2char numAl olga pam plotLine plotSeg plotTree pps \
randomizeSeq ranDot ranseq rep2plot \
repeater revComp rpois sass sblast sequencer shustring simNorm simOrf sops \
testMeans translate travTree upgma var watterson wrapSeq

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
tangle:
	test -d bin || mkdir bin
	for pack in $(packs); do \
		make tangle -C $$pack; \
	done
	printf "" > progs.txt
	for prog in $(progs); do \
		make tangle -C $$prog; \
		echo $$prog >> progs.txt; \
		cp $$prog/$$prog bin; \
	done
.PHONY: weave
weave:
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
