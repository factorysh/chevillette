build: bin
	go build -o bin/chevillette cli/main.go

test:
	go test -cover github.com/factorysh/chevillette/log

bin:
	mkdir -p bin
