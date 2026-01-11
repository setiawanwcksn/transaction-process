APP=flip
PKG=./...
BIN=bin/$(APP)

build:
	CGO_ENABLED=0 go build -o bin/flip cmd/server/main.go

run: build
	./bin/flip

test:
	GOFLAGS="-ldflags=-linkmode=external" go test $(PKG)

testv:
	go test -v $(PKG)

race:
	go test -race $(PKG)

cover:	
	GOFLAGS="-ldflags=-linkmode=external" go test -coverprofile=coverage.out $(PKG)
	go tool cover -func=coverage.out

clean:
	rm -rf $(BIN) coverage.out

lint:
	golangci-lint run

.PHONY: run build test testv race cover clean lint
