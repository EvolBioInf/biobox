all : newick

newick: newick.go
	go build newick.go
newick.go: newick.org
	awk -f ../scripts/preTangle.awk newick.org | bash ../scripts/org2nw | notangle -Rnewick.go | gofmt > newick.go
test: newick_test.go newick.go
	go test -v
newick_test.go: newick.org
	bash ../scripts/org2nw newick.org | notangle -Rnewick_test.go | gofmt > newick_test.go

clean:
	rm -f *.go
