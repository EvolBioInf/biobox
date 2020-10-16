all: bioboxDoc.tex cchar.tex getSeq.tex util.tex wrapSeq.tex
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
util.tex: ../util/util.org
	bash ../scripts/org2nw ../util/util.org | awk -f ../scripts/preWeave.awk | noweave -n -x > util.tex
wrapSeq.tex: ../wrapSeq/wrapSeq.org
	bash ../scripts/org2nw ../wrapSeq/wrapSeq.org | noweave -n -x | sed 's/_/\\_/g' > wrapSeq.tex
clean:
	rm -f cchar.tex getSeq.tex util.tex wrapSeq.tex *.pdf *.aux *.bbl *.blg *.dvi *.log *.toc 
