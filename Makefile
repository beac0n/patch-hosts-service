dev:
	go get github.com/oxequa/realize; realize start
build:
	go build -o build/patchbay-server src/main.go
clean:
	rm -rf build
	rm -f coverage.html
run:
	build/patchbay-server
cov:
	go test -coverpkg=./... -cover -coverprofile coverage.html -v ./...
	go tool cover -html=coverage.html
test:
	go test -v ./...