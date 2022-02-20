build: bin
	go build -o bin/chevillette cli/main.go

test:
	go test -cover \
		github.com/factorysh/chevillette/log \
		github.com/factorysh/chevillette/memory \
		github.com/factorysh/chevillette/pattern

bin:
	mkdir -p bin
