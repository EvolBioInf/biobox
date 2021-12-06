src = al.tex blast2dot.tex cchar.tex cutSeq.tex dnaDist.tex drawKt.tex drawSt.tex fasta2tab.tex \
	genTree.tex getSeq.tex histogram.tex keyMat.tex mum2plot.tex mutator.tex \
	naiveMatcher.tex nj.tex numAl.tex pam.tex plotLine.tex plotSeg.tex plotTree.tex \
	ranseq.tex randomizeSeq.tex ranDot.tex 	rep2plot.tex rpois.tex sblast.tex shustring.tex \
	simNorm.tex simOrf.tex testMeans.tex translate.tex travTree.tex repeater.tex revComp.tex \
	upgma.tex util.tex var.tex watterson.tex wrapSeq.tex

all: $(src) bioboxDoc.pdf 
%.tex:	../*/%.org bioboxDoc.tex
	bash ../scripts/org2nw $< | noweave -n -x | sed 's/_test/\\_test/g' > $@
bioboxDoc.pdf: $(src) rd1.ps rd2.ps
	latex bioboxDoc
	bibtex bioboxDoc
	latex bioboxDoc
	latex bioboxDoc
	latex bioboxDoc
	dvips bioboxDoc -o -q
	ps2pdf bioboxDoc.ps
rd1.ps: ../ranDot/rd1.dot
	neato -T ps ../ranDot/rd1.dot > rd1.ps
rd2.ps: ../ranDot/rd2.dot
	neato -T ps ../ranDot/rd2.dot > rd2.ps
clean:
	rm -f $(src) *.pdf *.aux *.bbl *.blg *.dvi *.log *.to rd1.ps rd2.ps
