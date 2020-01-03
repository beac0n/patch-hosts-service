setup:
	go get github.com/oxequa/realize
dev:
	realize start
build:
	go build -race -o build/patchbay-server src/main.go
clean:
	rm -rf build
run:
	build/patchbay-server
