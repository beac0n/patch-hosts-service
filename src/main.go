package main

import (
	"./handlers/pubsub"
	"flag"
	"log"
	"net/http"
)

type RequestHandler struct {
	pubSubRequestHandler *pubsub.RequestHandler
}

func (requestHandler *RequestHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	requestHandler.pubSubRequestHandler.ServeHttp(responseWriter, request)
}

func main() {
	host := flag.String("host", "0.0.0.0:9001", "host and port where this application should run")
	maxReqSizeInMb := flag.Int64("max_req_size", 10, "maximum request size in MB")

	flag.Parse()

	log.Println("running on", *host)

	requestHandler := &RequestHandler{
		pubSubRequestHandler: pubsub.NewRequestHandler(*maxReqSizeInMb * 1000 * 1000),
	}

	if err := http.ListenAndServe(*host, requestHandler); err != nil {
		log.Fatal("FATAL ERROR:", err)
	}
}
