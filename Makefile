.PHONY: all
all: bin/purl

go.mod go.sum:
	go mod tidy

bin/purl: main.go cli/*.go
	go build -o bin/purl main.go

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck -checks="all,-ST1000" ./...

.PHONY: clean
clean:
	rm -rf bin/*
