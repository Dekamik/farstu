all: tidy build

tidy:
	go mod tidy

build:
	go build -o bin/farstu cmd/main.go

install:
	cp bin/farstu /usr/local/bin

run:
	go run cmd/main.go

clean:
	rm bin/farstu
	rm internal/views/*_templ.go
