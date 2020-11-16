all: al.tex bioboxDoc.tex cchar.tex cutSeq.tex getSeq.tex naiveMatcher.tex numAl.tex ranseq.tex randomizeSeq revComp.tex util.tex var.tex watterson.tex wrapSeq.tex
	latex bioboxDoc
	bibtex bioboxDoc
	latex bioboxDoc
	latex bioboxDoc
	latex bioboxDoc
	dvipdf bioboxDoc
al.tex: ../al/al.org
	bash ../scripts/org2nw ../al/al.org   | noweave -n -x | sed 's/_test/\\_test/g' > al.tex
cchar.tex: ../cchar/cchar.org
	bash ../scripts/org2nw ../cchar/cchar.org   | noweave -n -x | sed 's/_test/\\_test/g' > cchar.tex
cutSeq.tex: ../cutSeq/cutSeq.org
	bash ../scripts/org2nw ../cutSeq/cutSeq.org   | noweave -n -x | sed 's/_test/\\_test/g' > cutSeq.tex
getSeq.tex: ../getSeq/getSeq.org
	bash ../scripts/org2nw ../getSeq/getSeq.org | noweave -n -x | sed 's/_test/\\_test/g' > getSeq.tex
naiveMatcher.tex: ../naiveMatcher/naiveMatcher.org
	bash ../scripts/org2nw ../naiveMatcher/naiveMatcher.org | noweave -n -x | sed 's/_test/\\_test/g' > naiveMatcher.tex
numAl.tex: ../numAl/numAl.org
	bash ../scripts/org2nw ../numAl/numAl.org | noweave -n -x | sed 's/_test/\\_test/g' > numAl.tex
randomizeSeq: ../randomizeSeq/randomizeSeq.org
	bash ../scripts/org2nw ../randomizeSeq/randomizeSeq.org | noweave -n -x | sed 's/_test/\\_test/g' > randomizeSeq.tex
ranseq.tex: ../ranseq/ranseq.org
	bash ../scripts/org2nw ../ranseq/ranseq.org | noweave -n -x | sed 's/_test/\\_test/g' > ranseq.tex	
revComp.tex: ../revComp/revComp.org
	bash ../scripts/org2nw ../revComp/revComp.org | noweave -n -x | sed 's/_test/\\_test/g' > revComp.tex	
util.tex: ../util/util.org
	bash ../scripts/org2nw ../util/util.org | awk -f ../scripts/preWeave.awk | noweave -n -x | sed 's/_test/\\_test/g' > util.tex
var.tex: ../var/var.org
	bash ../scripts/org2nw ../var/var.org | noweave -n -x | sed 's/_test/\\_test/g' > var.tex
watterson.tex: ../watterson/watterson.org
	bash ../scripts/org2nw ../watterson/watterson.org | noweave -n -x | sed 's/_test/\\_test/g' > watterson.tex
wrapSeq.tex: ../wrapSeq/wrapSeq.org
	bash ../scripts/org2nw ../wrapSeq/wrapSeq.org | noweave -n -x | sed 's/_test/\\_test/g' > wrapSeq.tex
clean:
	rm -f al.tex cchar.tex cutSeq.tex getSeq.tex numAl.tex util.tex var.tex wrapSeq.tex *.pdf *.aux *.bbl *.blg *.dvi *.log *.toc 
