version = 0.1.0-nightly

all: build run

build: front backend


backend:
	go build -ldflags "-s -w"

.PHONY: front
.ONESHELL: front
front:
	cd front
	pnpm build
clean:
	rm plakken
	rm -r build/

lint:
	golangci-lint run

test:
	go test ./...

run:
	./status
