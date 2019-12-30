package main

import (
	"./handlers"
	"./utils"
	"bytes"
	"log"
	"net/http"
	"sync"
)

var address = ":9001"

type Channel struct {
	data     map[string]chan bytes.Buffer
	com      map[string]chan struct{}
	mux      *sync.Mutex
	chanSize uint
}

var pubSubChannels = &Channel{data: make(map[string]chan bytes.Buffer), com: make(map[string]chan struct{}), mux: &sync.Mutex{}, chanSize: 100}
var defaultChannels = &Channel{data: make(map[string]chan bytes.Buffer), com: make(map[string]chan struct{}), mux: &sync.Mutex{}}

func (channels *Channel) getChannels(path string) (chan bytes.Buffer, chan struct{}) {
	channels.mux.Lock()
	defer channels.mux.Unlock()

	_, dataOk := channels.data[path]
	if !dataOk {
		channels.data[path] = make(chan bytes.Buffer, channels.chanSize)
	}

	_, comOk := channels.com[path]
	if !comOk {
		channels.com[path] = make(chan struct{}, channels.chanSize)
	}

	return channels.data[path], channels.com[path]
}

func requestHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodGet) && (request.Method != http.MethodPost) {
		responseWriter.WriteHeader(http.StatusBadRequest)
		_, err := responseWriter.Write([]byte("wrong http method"))
		utils.LogError(err, request)
		return
	}

	var dataChannel chan bytes.Buffer
	var comChannel chan struct{}
	var handleProducer func(dataChannel chan bytes.Buffer, comChannel chan struct{}, request *http.Request, responseWriter http.ResponseWriter)
	var handleConsumer func(dataChannel chan bytes.Buffer, comChannel chan struct{}, request *http.Request, responseWriter http.ResponseWriter)

	pubSubKeys, ok := request.URL.Query()["pubsub"]
	if ok && len(pubSubKeys) == 1 && pubSubKeys[0] == "true" {
		dataChannel, comChannel = pubSubChannels.getChannels(request.URL.Path)
		handleProducer = handlers.HandleProducerPubSub
		handleConsumer = handlers.HandleConsumerPubSub
	} else {
		dataChannel, comChannel = defaultChannels.getChannels(request.URL.Path)
		handleProducer = handlers.HandleProducerDefault
		handleConsumer = handlers.HandleConsumerDefault
	}

	if request.Method == http.MethodPost {
		handleProducer(dataChannel, comChannel, request, responseWriter)
	} else if request.Method == http.MethodGet {
		handleConsumer(dataChannel, comChannel, request, responseWriter)
	}
}

func main() {
	log.Println("running on", address)

	if err := http.ListenAndServe(address, http.HandlerFunc(requestHandler)); err != nil {
		log.Fatal(err)
	}
}
