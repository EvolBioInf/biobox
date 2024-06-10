VERSION = $(shell bash ../scripts/getVersion.sh)
DATE = $(shell bash ../scripts/getDate.sh)

EXE = ranseq
VF = -X github.com/evolbioinf/biobox/util.version=$(VERSION)
DF = -X github.com/evolbioinf/biobox/util.date=$(DATE)
BUILD = go build -ldflags "$(VF) $(DF)" $(EXE).go
NW = $(shell which noweb)

$(EXE): $(EXE).go
	$(BUILD)
tangle: $(EXE).go $(EXE)_test.go
$(EXE).go: $(EXE).org
	if [ $(NW) != "" ]; then\
		bash ../scripts/org2nw $(EXE).org | notangle -R$(EXE).go | gofmt > $(EXE).go;\
	fi
test: $(EXE) $(EXE)_test.go
	go test -v
$(EXE)_test.go: $(EXE).org
	if [ $(NW) != "" ]; then\
		bash ../scripts/org2nw $(EXE).org | notangle -R$(EXE)_test.go | gofmt > $(EXE)_test.go;\
	fi
clean:
	rm -f $(EXE) *.go
