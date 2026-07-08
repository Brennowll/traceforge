.PHONY: test build run lint race clean

test:
	go test ./...

build:
	go build -o bin/traceforge ./cmd/traceforge

run:
	go run ./cmd/traceforge

lint:
	go vet ./...

race:
	go test -race ./...

clean:
	rm -rf bin report.html
