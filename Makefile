progs = al blast2dot cchar cutSeq dnaDist drawKt drawSt fasta2tab getSeq keyMat mum2plot mutator naiveMatcher numAl pam \
plotLine plotSeg randomizeSeq ranDot ranseq rep2plot repeater revComp rpois shustring simNorm testMeans watterson wrapSeq

all:
	for prog in $(progs); do \
		make -C $$prog; \
		cp $$prog/$$prog bin; \
	done
.PHONY: doc
doc:
	make -C doc
clean:
	for prog in $(progs) doc; do \
		make clean -C $$prog; \
	done
test:
	for prog in $(progs); do \
		make test -C $$prog; \
	done
