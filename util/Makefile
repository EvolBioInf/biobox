all : util

util: util.go
	go build util.go
util.go: util.org
	awk -f ../scripts/preTangle.awk util.org | bash ../scripts/org2nw | notangle -Rutil.go | gofmt > util.go

.PHONY:	doc
doc:
	make -C ../doc

clean:
	rm -f *.go
	make clean -C ../doc
