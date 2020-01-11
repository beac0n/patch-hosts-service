package main

import (
	"./handlers/mpmc"
	"./handlers/pubsub"
	"./handlers/reqres"
	"flag"
	"log"
	"net/http"
)

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
