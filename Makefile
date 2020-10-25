export VERSION = $(shell bash ./scripts/version.sh) # The single, trailing blank is essential
export DATE    = $(shell bash ./scripts/date.sh) #    The single, trailing blank is essential

all:
	make -C util
	make -C cchar
	cp cchar/cchar bin
	make -C cutSeq
	cp cutSeq/cutSeq bin
	make -C getSeq
	cp getSeq/getSeq bin
	make -C randomizeSeq
	cp randomizeSeq/randomizeSeq bin
	make -C ranseq
	cp ranseq/ranseq bin
	make -C revComp
	cp revComp/revComp bin
	make -C wrapSeq
	cp wrapSeq/wrapSeq bin
	make -C var
	cp var/var bin
.PHONY: doc
doc:
	make -C doc
clean:
	make clean -C cchar
	make clean -C cutSeq
	make clean -C getSeq
	make clean -C doc
	make clean -C randomizeSeq
	make clean -C ranseq
	make clean -C revComp
	make clean -C util
	make clean -C wrapSeq
	make clean -C var
test:
	make test -C cchar
	make tets -C cutSeq
	make test -C getSeq
	make test -C randomizeSeq
	make test -C ranseq
	make test -C revComp
	make test -C wrapSeq
	make test -C var
