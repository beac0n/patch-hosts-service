dev:
	go get github.com/oxequa/realize; realize start
build:
	go build -o build/patch-hosts-service-linux-amd64 src/main.go
clean:
	rm -rf build
	rm -f coverage.html
run:
	build/patch-hosts-service-linux-amd64
cov:
	go test -tags=test -coverpkg=./... -cover -coverprofile coverage.html -v ./...
	go tool cover -html=coverage.html
test:
	go test -tags=test -v ./...
docker_build:
	docker build -t patch-hosts-service .
docker_run:
	docker run -p 9001:9001 -it patch-hosts-service
docker: docker_build docker_run