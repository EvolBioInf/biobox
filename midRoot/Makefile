VERSION = $(shell git describe)
DATE := $(shell git log | grep -m 1 Date | sed -r 's/Date: +[A-Z][a-z]+ ([A-Z][a-z]+) ([0-9]+) [^ ]+ ([0-9]+) .+/\2_\1_\3/')

EXE = midRoot
VF = -X github.com/evolbioinf/biobox/util.version=$(VERSION)
DF = -X github.com/evolbioinf/biobox/util.date=$(DATE)
BUILD = go build -ldflags "$(VF) $(DF)" $(EXE).go

$(EXE):
	$(BUILD)
tangle: $(EXE).go $(EXE)_test.go
$(EXE).go: $(EXE).org
	bash ../scripts/org2nw $(EXE).org | notangle -R$(EXE).go | gofmt > $(EXE).go
test: $(EXE) $(EXE)_test.go
	go test -v
$(EXE)_test.go: $(EXE).org
	bash ../scripts/org2nw $(EXE).org | notangle -R$(EXE)_test.go | gofmt > $(EXE)_test.go
clean:
	rm -f $(EXE) *.go
