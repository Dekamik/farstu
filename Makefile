all: templ tidy build

deps:
	go install github.com/a-h/templ/cmd/templ@latest

tidy:
	go mod tidy

templ:
	templ generate

build:
	go build -o bin/farstu cmd/main.go
