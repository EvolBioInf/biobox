all:
	make -C util
	make -C al
	cp al/al bin
	make -C cchar
	cp cchar/cchar bin
	make -C cutSeq
	cp cutSeq/cutSeq bin
	make -C dnaDist
	cp dnaDist/dnaDist bin
	make -C drawKt
	cp drawKt/drawKt bin
	make -C drawSt
	cp drawSt/drawSt bin
	make -C fasta2tab
	cp fasta2tab/fasta2tab bin
	make -C getSeq
	cp getSeq/getSeq bin
	make -C keyMat
	cp keyMat/keyMat bin
	make -C mutator
	cp mutator/mutator bin
	make -C naiveMatcher
	cp naiveMatcher/naiveMatcher bin
	make -C numAl
	cp numAl/numAl bin
	make -C pam
	cp pam/pam bin
	make -C randomizeSeq
	cp randomizeSeq/randomizeSeq bin
	make -C ranseq
	cp ranseq/ranseq bin
	make -C repeater
	cp repeater/repeater bin
	make -C revComp
	cp revComp/revComp bin
	make -C shustring
	cp shustring/shustring bin
	make -C simNorm
	cp simNorm/simNorm bin
	make -C testMeans
	cp testMeans/testMeans bin
	make -C watterson
	cp watterson/watterson bin
	make -C wrapSeq
	cp wrapSeq/wrapSeq bin
	make -C var
	cp var/var bin
.PHONY: doc
doc:
	make -C doc
clean:
	make clean -C al
	make clean -C cchar
	make clean -C cutSeq
	make clean -C dnaDist
	make clean -C doc
	make clean -C drawKt
	make clean -C drawSt
	make clean -C fasta2tab
	make clean -C getSeq
	make clean -C keyMat
	make clean -C mutator
	make clean -C naiveMatcher
	make clean -C numAl
	make clean -C pam
	make clean -C randomizeSeq
	make clean -C ranseq
	make clean -C repeater
	make clean -C revComp
	make clean -C shustring
	make clean -C simNorm
	make clean -C testMeans
	make clean -C util
	make clean -C var
	make clean -C watterson
	make clean -C wrapSeq
test:
	make test -C al
	make test -C cchar
	make test -C cutSeq
	make test -C dnaDist
	make test -C drawKt
	make test -C drawSt
	make test -C fasta2tab
	make test -C getSeq
	make test -C keyMat
	make test -C mutator
	make test -C naiveMatcher
	make test -C numAl
	make test -C pam
	make test -C randomizeSeq
	make test -C ranseq
	make test -C repeater
	make test -C revComp
	make test -C shustring
	make test -C simNorm
	make test -C testMeans
	make test -C util
	make test -C var
	make test -C watterson
	make test -C wrapSeq
