package main

import (
	"./handlers/mpmc"
	"./handlers/pubsub"
	"flag"
	"log"
	"net/http"
	"strings"
)

type RequestHandler struct {
	pubSubRequestHandler *pubsub.RequestHandler
	mpmcRequestHandler   *mpmc.RequestHandler
}

func (requestHandler *RequestHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if isCorrectPath(request, "/pubsub") {
		requestHandler.pubSubRequestHandler.ServeHttp(responseWriter, request)
		return
	}

	if isCorrectPath(request, "/queue") {
		requestHandler.mpmcRequestHandler.ServeHttp(responseWriter, request)
		return
	}

	http.Error(responseWriter, "", http.StatusNotFound)
}

func isCorrectPath(request *http.Request, path string) bool {
	return strings.HasPrefix(request.URL.Path, path) && request.URL.Path != path
}

func main() {
	host := flag.String("host", "0.0.0.0:9001", "host and port where this application should run")
	maxReqSizeInMb := flag.Int64("max_req_size", 10, "maximum request size in MB")

	flag.Parse()

	log.Println("running on", *host)

	maxReqSize := *maxReqSizeInMb * 1000 * 1000
	requestHandler := &RequestHandler{
		pubSubRequestHandler: pubsub.NewRequestHandler(maxReqSize),
		mpmcRequestHandler:   mpmc.NewRequestHandler(maxReqSize),
	}

	if err := http.ListenAndServe(*host, requestHandler); err != nil {
		log.Fatal("FATAL ERROR:", err)
	}
}
