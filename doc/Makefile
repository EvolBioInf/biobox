all: al.tex bioboxDoc.tex cchar.tex cutSeq.tex dnaDist.tex drawKt.tex getSeq.tex mutator.tex naiveMatcher.tex numAl.tex pam.tex ranseq.tex randomizeSeq.tex simNorm.tex testMeans.tex revComp.tex util.tex var.tex watterson.tex wrapSeq.tex
	latex bioboxDoc
	bibtex bioboxDoc
	latex bioboxDoc
	latex bioboxDoc
	latex bioboxDoc
	dvips bioboxDoc -o -q
	ps2pdf bioboxDoc.ps
al.tex: ../al/al.org
	bash ../scripts/org2nw ../al/al.org   | noweave -n -x | sed 's/_test/\\_test/g' > al.tex
cchar.tex: ../cchar/cchar.org
	bash ../scripts/org2nw ../cchar/cchar.org   | noweave -n -x | sed 's/_test/\\_test/g' > cchar.tex
cutSeq.tex: ../cutSeq/cutSeq.org
	bash ../scripts/org2nw ../cutSeq/cutSeq.org   | noweave -n -x | sed 's/_test/\\_test/g' > cutSeq.tex
dnaDist.tex: ../dnaDist/dnaDist.org
	bash ../scripts/org2nw ../dnaDist/dnaDist.org   | noweave -n -x | sed 's/_test/\\_test/g' > dnaDist.tex
drawKt.tex: ../drawKt/drawKt.org
	bash ../scripts/org2nw ../drawKt/drawKt.org   | noweave -n -x | sed 's/_test/\\_test/g' > drawKt.tex
getSeq.tex: ../getSeq/getSeq.org
	bash ../scripts/org2nw ../getSeq/getSeq.org | noweave -n -x | sed 's/_test/\\_test/g' > getSeq.tex
mutator.tex: ../mutator/mutator.org
	bash ../scripts/org2nw ../mutator/mutator.org | noweave -n -x | sed 's/_test/\\_test/g' > mutator.tex
naiveMatcher.tex: ../naiveMatcher/naiveMatcher.org
	bash ../scripts/org2nw ../naiveMatcher/naiveMatcher.org | noweave -n -x | sed 's/_test/\\_test/g' > naiveMatcher.tex
numAl.tex: ../numAl/numAl.org
	bash ../scripts/org2nw ../numAl/numAl.org | noweave -n -x | sed 's/_test/\\_test/g' > numAl.tex
pam.tex: ../pam/pam.org
	bash ../scripts/org2nw ../pam/pam.org | noweave -n -x | sed 's/_test/\\_test/g' > pam.tex
randomizeSeq.tex: ../randomizeSeq/randomizeSeq.org
	bash ../scripts/org2nw ../randomizeSeq/randomizeSeq.org | noweave -n -x | sed 's/_test/\\_test/g' > randomizeSeq.tex
ranseq.tex: ../ranseq/ranseq.org
	bash ../scripts/org2nw ../ranseq/ranseq.org | noweave -n -x | sed 's/_test/\\_test/g' > ranseq.tex	
revComp.tex: ../revComp/revComp.org
	bash ../scripts/org2nw ../revComp/revComp.org | noweave -n -x | sed 's/_test/\\_test/g' > revComp.tex
simNorm.tex: ../simNorm/simNorm.org
	bash ../scripts/org2nw ../simNorm/simNorm.org | noweave -n -x | sed 's/_test/\\_test/g' > simNorm.tex
testMeans.tex: ../testMeans/testMeans.org
	bash ../scripts/org2nw ../testMeans/testMeans.org | noweave -n -x | sed 's/_test/\\_test/g' > testMeans.tex	
util.tex: ../util/util.org
	bash ../scripts/org2nw ../util/util.org | awk -f ../scripts/preWeave.awk | noweave -n -x | sed 's/_test/\\_test/g' > util.tex
var.tex: ../var/var.org
	bash ../scripts/org2nw ../var/var.org | noweave -n -x | sed 's/_test/\\_test/g' > var.tex
watterson.tex: ../watterson/watterson.org
	bash ../scripts/org2nw ../watterson/watterson.org | noweave -n -x | sed 's/_test/\\_test/g' > watterson.tex
wrapSeq.tex: ../wrapSeq/wrapSeq.org
	bash ../scripts/org2nw ../wrapSeq/wrapSeq.org | noweave -n -x | sed 's/_test/\\_test/g' > wrapSeq.tex
clean:
	rm -f al.tex cchar.tex cutSeq.tex dnaDist.tex getSeq.tex mutator.tex naiveMatcher.tex numAl.tex pam.tex ranseq.tex randomizeSeq.tex simNorm.tex testMeans.tex util.tex var.tex watterson.tex wrapSeq.tex *.pdf *.aux *.bbl *.blg *.dvi *.log *.toc 
