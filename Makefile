all: build

build:
	go build -o ./bin/vote ./cmd/vote

test:
	go test -race -cover -count=1 ./...