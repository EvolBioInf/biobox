progs = al blast2dot bwt cres cutSeq dnaDist drag drawf drawGenes drawKt drawSt fasta2tab geco genTree getSeq histogram kerror keyMat midRoot maf mtf \
mum2plot mutator naiveMatcher nj numAl pam plotLine plotSeg plotTree pps randomizeSeq ranDot ranseq rep2plot \
repeater revComp rpois sblast sequencer shustring simNorm simOrf testMeans translate travTree upgma var watterson \
wrapSeq
packs = util

all:
	test -d bin || mkdir bin
	for pack in $(packs); do \
		make -C $$pack; \
	done
	for prog in $(progs); do \
		make -C $$prog; \
		cp $$prog/$$prog bin; \
	done
.PHONY: doc
doc:
	make -C doc
clean:
	for prog in $(progs) $(packs) doc; do \
		make clean -C $$prog; \
	done
test:
	for prog in $(progs) $(packs); do \
		make test -C $$prog; \
	done
init:
	for pack in $(packs); do \
		cd $$pack; go mod init; cd ..; \
	done
	for prog in $(progs); do \
		cd $$prog; go mod init; cd ..; \
	done
tidy:
	for pack in $(packs); do \
		cd $$pack; go mod tidy; cd ..; \
	done
	for prog in $(progs); do \
		cd $$prog; go mod tidy; cd ..; \
	done
