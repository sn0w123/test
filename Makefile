build:
	go build -o bin/go-scaffold ./cmd/go-scaffold

install:
	go install ./cmd/go-scaffold

run: build
	./bin/go-scaffold new order