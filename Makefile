build: bin
	rm -f bien/chevillette
	go build -o ./bin/chevillette cli/main.go

build-linux:
	make build GOOS=linux
	if [ "upx not found" != "$(shell which upx)" ]; then upx bin/chevillette; fi

build-with-docker: bin
	mkdir -p .cache
	docker run -ti --rm \
		-v `pwd`:/src/ \
		-v `pwd`/.cache:/.cache \
		-w /src \
		-u `id -u` \
		golang:1.17-bullseye \
		make build
	file bin/chevillette

test:
	go test -cover \
		github.com/factorysh/chevillette/log \
		github.com/factorysh/chevillette/memory \
		github.com/factorysh/chevillette/pattern

bin:
	mkdir -p bin

clean:
	rm -rf bin
	#rm -rf .cache
