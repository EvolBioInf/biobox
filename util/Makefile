NW = $(shell which noweb)

all : util

util: util.go
	go build util.go
tangle: util.go util_test.go
	go build util.go
util.go: util.org
	if [ $(NW) != "" ]; then\
		awk -f ../scripts/preTangle.awk util.org | bash ../scripts/org2nw | notangle -Rutil.go | gofmt > util.go;\
	fi
test: util_test.go util.go
	go test -v
util_test.go: util.org
	if [ $(NW) != "" ]; then\
		bash ../scripts/org2nw util.org | notangle -Rutil_test.go | gofmt > util_test.go;\
	fi

.PHONY:	doc
doc:
	make -C ../doc

clean:
	rm -f *.go
	make clean -C ../doc
