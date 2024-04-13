src = al.tex blast2dot.tex bwt.tex clac.tex cres.tex cutSeq.tex dnaDist.tex \
drag.tex drawf.tex drawGenes.tex drawKt.tex drawSt.tex fasta2tab.tex \
geco.tex genTree.tex getSeq.tex histogram.tex huff.tex hut.tex \
kerror.tex keyMat.tex maf.tex midRoot.tex mtf.tex mum2plot.tex mutator.tex \
naiveMatcher.tex nj.tex num2char.tex numAl.tex olga.tex pam.tex plotLine.tex \
plotSeg.tex plotTree.tex pps.tex ranseq.tex randomizeSeq.tex ranDot.tex \
rep2plot.tex rpois.tex sass.tex sblast.tex sequencer.tex shustring.tex \
simNorm.tex simOrf.tex sops.tex testMeans.tex translate.tex travTree.tex \
repeater.tex revComp.tex upgma.tex util.tex var.tex watterson.tex wrapSeq.tex

date = $(shell git log | grep Date | head -n 1 | sed -r 's/Date: +[A-Z][a-z]+ ([A-Z][a-z]+) ([0-9]+) [^ ]+ ([0-9]+) .+/\2_\1_\3/')
version = $(shell git describe)

all: $(src) biobox.pdf 
%.tex:	../*/%.org biobox.tex
	bash ../scripts/org2nw $< | noweave -n -x | sed 's/_test/\\_test/g' > $@
biobox.pdf: $(src) rd1.ps rd2.ps
	echo $(date) | tr '_' ' ' > date.txt
	echo $(version) | tr '-' ' ' | awk '{printf "%s", $$1; if ($$2) printf "-%s", $$2; printf "\n"}' > version.txt
	latex biobox
	bibtex biobox
	latex biobox
	latex biobox
	latex biobox
	dvips biobox -o -q
	ps2pdf -dALLOWPSTRANSPARENCY biobox.ps
	if [ -d /mnt/c/Users/haubold/ownCloud\ -\ haubold@evolbio.mpg.de@owncloud.gwdg.de/docs/ ]; then \
		cp biobox.pdf /mnt/c/Users/haubold/ownCloud\ -\ haubold@evolbio.mpg.de@owncloud.gwdg.de/docs/; \
	fi
rd1.ps: ../ranDot/rd1.dot
	neato -T ps ../ranDot/rd1.dot > rd1.ps
rd2.ps: ../ranDot/rd2.dot
	neato -T ps ../ranDot/rd2.dot > rd2.ps
clean:
	rm -f $(src) *.pdf *.aux *.bbl *.blg *.dvi *.log *.to *.toc rd1.ps rd2.ps
