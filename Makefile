all: tidy templ build

tidy:
	go mod tidy

templ:
	templ generate

build:
	go build -o bin/farstu cmd/main.go
