export VERSION = $(shell bash ../scripts/version.sh) # The single, trailing blank is essential
export DATE    = $(shell bash ../scripts/date.sh) #    The single, trailing blank is essential

EXE = getSeq

all: $(EXE)

$(EXE): $(EXE).go
	go build -ldflags "-X github.com/evolbioinf/biobox/util.version=$(VERSION) -X github.com/evolbioinf/biobox/util.date=$(DATE)" $(EXE).go
$(EXE).go: $(EXE).org
	bash ../scripts/org2nw $(EXE).org | notangle -R$(EXE).go | gofmt > $(EXE).go
test: $(EXE) $(EXE)_test.go
	go test -v
$(EXE)_test.go: $(EXE).org
	bash ../scripts/org2nw $(EXE).org | notangle -R$(EXE)_test.go | gofmt > $(EXE)_test.go
clean:
	rm -f $(EXE) *.go
