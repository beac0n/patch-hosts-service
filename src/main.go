package main

import (
	"./handlers"
	"flag"
	"log"
	"net/http"
)

type RequestHandler struct {
	pubSubRequestHandler *handlers.PubSubRequestHandler
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
		pubSubRequestHandler: handlers.NewPubSubRequestHandler(*maxReqSizeInMb),
	}

	if err := http.ListenAndServe(*host, requestHandler); err != nil {
		log.Fatal("FATAL ERROR:", err)
	}
}
