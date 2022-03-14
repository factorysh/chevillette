build: bin
	go build -o bin/chevillette cli/main.go

build-linux:
	make build GOOS=linux
	upx bin/chevillette

test:
	go test -cover \
		github.com/factorysh/chevillette/log \
		github.com/factorysh/chevillette/memory \
		github.com/factorysh/chevillette/pattern

bin:
	mkdir -p bin
