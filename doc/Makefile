all: bioboxDoc.tex cchar.tex getSeq.tex ranseq.tex util.tex var.tex wrapSeq.tex
	latex bioboxDoc
	bibtex bioboxDoc
	latex bioboxDoc
	latex bioboxDoc
	latex bioboxDoc
	dvipdf bioboxDoc
cchar.tex: ../cchar/cchar.org
	bash ../scripts/org2nw ../cchar/cchar.org   | noweave -n -x | sed 's/_/\\_/g' > cchar.tex
getSeq.tex: ../getSeq/getSeq.org
	bash ../scripts/org2nw ../getSeq/getSeq.org | noweave -n -x | sed 's/_/\\_/g' > getSeq.tex	
ranseq.tex: ../ranseq/ranseq.org
	bash ../scripts/org2nw ../ranseq/ranseq.org | noweave -n -x | sed 's/_/\\_/g' > ranseq.tex	
util.tex: ../util/util.org
	bash ../scripts/org2nw ../util/util.org | awk -f ../scripts/preWeave.awk | noweave -n -x > util.tex
var.tex: ../var/var.org
	bash ../scripts/org2nw ../var/var.org | noweave -n -x | sed 's/_/\\_/g' > var.tex
wrapSeq.tex: ../wrapSeq/wrapSeq.org
	bash ../scripts/org2nw ../wrapSeq/wrapSeq.org | noweave -n -x | sed 's/_/\\_/g' > wrapSeq.tex
clean:
	rm -f cchar.tex getSeq.tex util.tex var.tex wrapSeq.tex *.pdf *.aux *.bbl *.blg *.dvi *.log *.toc 
