export VERSION = $(shell bash ./scripts/version.sh) # The single, trailing blank is essential
export DATE    = $(shell bash ./scripts/date.sh) #    The single, trailing blank is essential

all:
	make -C util
	make -C cchar
	make -C getSeq
.PHONY: doc
doc:
	make -C doc

clean:
	make clean -C cchar
	make clean -C getSeq
	make clean -C doc
	make clean -C util
