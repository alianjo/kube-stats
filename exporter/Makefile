all: build

build:
	go build -o bin/main main.go

run:
	go run main.go

test:
	go test -v ./...

docker-build:
	docker build -t myapp .

docker-run:
	docker run -p 8080:8080 myapp

clean:
	rm -f bin/main

.PHONY: all build run test docker-build docker-run clean