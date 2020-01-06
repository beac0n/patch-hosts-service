package main

import (
	"./handlers/mpmc"
	"./handlers/pubsub"
	"flag"
	"log"
	"net/http"
	"strings"
)

type ReqHandler struct {
	pubSubReqHandler *pubsub.ReqHandler
	mpmcReqHandler   *mpmc.ReqHandler
}

func (reqHandler *ReqHandler) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	if isCorrectPath(req, "/pubsub") {
		reqHandler.pubSubReqHandler.ServeHTTP(resWriter, req)
		return
	}

	if isCorrectPath(req, "/queue") {
		reqHandler.mpmcReqHandler.ServeHTTP(resWriter, req)
		return
	}

	http.Error(resWriter, "", http.StatusNotFound)
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
	reqHandler := &ReqHandler{
		pubSubReqHandler: pubsub.NewReqHandler(maxReqSize),
		mpmcReqHandler:   mpmc.NewReqHandler(maxReqSize),
	}

	if err := http.ListenAndServe(*host, reqHandler); err != nil {
		log.Fatal("FATAL ERROR:", err)
	}
}
