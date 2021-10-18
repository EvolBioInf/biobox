VERSION = $(shell git describe)
DATE := $(shell git log | grep Date | head -n 1 | sed -r 's/Date: +[A-Z][a-z]+ ([A-Z][a-z]+) ([0-9]+) [^ ]+ ([0-9]+) .+/\2_\1_\3/')

EXE = plotLine

all: $(EXE)

$(EXE): $(EXE).go lineTmpl.txt
	go build -ldflags "-X github.com/evolbioinf/biobox/util.version=$(VERSION) -X github.com/evolbioinf/biobox/util.date=$(DATE)" $(EXE).go
$(EXE).go: $(EXE).org
	bash ../scripts/org2nw $(EXE).org | notangle -R$(EXE).go | gofmt > $(EXE).go
lineTmpl.txt: $(EXE).org
	bash ../scripts/org2nw $(EXE).org | notangle -RlineTmpl.txt > lineTmpl.txt
test: $(EXE) $(EXE)_test.go
	go test -v
$(EXE)_test.go: $(EXE).org
	bash ../scripts/org2nw $(EXE).org | notangle -R$(EXE)_test.go | gofmt > $(EXE)_test.go
clean:
	rm -f $(EXE) *.go template.txt
