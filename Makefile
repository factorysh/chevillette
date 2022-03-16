build: bin
	go build -o bin/chevillette cli/main.go

build-linux:
	make build GOOS=linux
	if [ "upx not found" != "$(shell which upx)" ]; then upx bin/chevillette; fi

test:
	go test -cover \
		github.com/factorysh/chevillette/log \
		github.com/factorysh/chevillette/memory \
		github.com/factorysh/chevillette/pattern

bin:
	mkdir -p bin
