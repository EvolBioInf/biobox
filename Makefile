export VERSION = $(shell bash ./scripts/version.sh) # The single, trailing blank is essential
export DATE    = $(shell bash ./scripts/date.sh) #    The single, trailing blank is essential

all : cchar

cchar: cchar.go
	go build -ldflags "-X main.VERSION=$(VERSION) main.DATE=$(DATE)" cchar.go
cchar.go: biobox.org
	bash scripts/org2nw biobox.org | notangle -Rcchar.go > cchar.go

clean:
	rm -f cchar *.go
