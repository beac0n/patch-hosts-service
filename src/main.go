package main

import (
	"./constants"
	"./handlers/mpmc"
	"./handlers/pubsub"
	"./handlers/reqres"
	"./utils"
	"flag"
	"log"
	"net/http"
)

type ReqHandler struct {
	pubSubReqHandler http.Handler
	mpmcReqHandler   http.Handler
	reqResReqHandler http.Handler
}

func (reqHandler *ReqHandler) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	if utils.IsCorrectPath(req, "/pubsub") {
		reqHandler.pubSubReqHandler.ServeHTTP(resWriter, req)
	} else if utils.IsCorrectPath(req, "/queue") {
		reqHandler.mpmcReqHandler.ServeHTTP(resWriter, req)
	} else if utils.IsCorrectPath(req, constants.Res) || utils.IsCorrectPath(req, constants.Req) {
		reqHandler.reqResReqHandler.ServeHTTP(resWriter, req)
	} else {
		http.Error(resWriter, "", http.StatusNotFound)
	}
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
		reqResReqHandler: reqres.NewReqHandler(maxReqSize),
	}

	if err := http.ListenAndServe(*host, reqHandler); err != nil {
		log.Fatal("FATAL ERROR:", err)
	}
}
