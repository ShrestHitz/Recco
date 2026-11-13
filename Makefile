.PHONY: run test build clean migrate-up

run:
	go run ./cmd/api/main.go

test:
	go test -v -race ./...

build:
	go build -o bin/api ./cmd/api/main.go

clean:
	rm -rf bin/
