src = al.tex cchar.tex cutSeq.tex dnaDist.tex drawKt.tex drawSt.tex fasta2tab.tex \
	getSeq.tex keyMat.tex mum2plot.tex mutator.tex naiveMatcher.tex numAl.tex \
	pam.tex ranseq.tex randomizeSeq.tex rpois.tex shustring.tex simNorm.tex simOrf.tex \
	testMeans.tex repeater.tex revComp.tex util.tex var.tex watterson.tex wrapSeq.tex

all: $(src) bioboxDoc.pdf
%.tex:	../*/%.org
	bash ../scripts/org2nw $< | noweave -n -x | sed 's/_test/\\_test/g' > $@
bioboxDoc.pdf: $(src)
	latex bioboxDoc
	bibtex bioboxDoc
	latex bioboxDoc
	latex bioboxDoc
	latex bioboxDoc
	dvips bioboxDoc -o -q
	ps2pdf bioboxDoc.ps
clean:
	rm -f $(src) *.pdf *.aux *.bbl *.blg *.dvi *.log *.to
