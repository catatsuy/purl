.PHONY: all
all: bin/purl

go.mod go.sum:
	go mod tidy

bin/purl: main.go cli/*.go go.mod go.sum
	go build -ldflags "-X github.com/catatsuy/purl/cli.Version=`git rev-list HEAD -n1`" -o bin/purl main.go

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck -checks="all,-ST1000" ./...

.PHONY: test
test:
	go test -cover -v ./...

.PHONY: clean
clean:
	rm -rf bin/*
