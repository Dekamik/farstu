all: templ tidy build

deps:
	go install github.com/a-h/templ/cmd/templ@latest

tidy:
	go mod tidy

templ:
	templ generate

build:
	go build -o bin/farstu cmd/main.go

install:
	cp bin/farstu /usr/local/bin

run:
	go run cmd/main.go

clean:
	rm bin/farstu
	rm internal/views/*_templ.go
