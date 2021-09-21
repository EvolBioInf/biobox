progs = al cchar cutSeq dnaDist drawKt drawSt fasta2tab getSeq keyMat mum2plot mutator naiveMatcher numAl pam \
randomizeSeq ranseq repeater revComp rpois shustring simNorm testMeans watterson wrapSeq

all:
	make -C util
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
