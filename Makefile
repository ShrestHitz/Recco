.PHONY: run test build clean web-install web-dev web-build

run:
	go run ./cmd/api/main.go

test:
	go test -v -race ./...

build:
	go build -o bin/api ./cmd/api/main.go

web-install:
	cd web && npm install

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build

clean:
	rm -rf bin/ web/dist/
