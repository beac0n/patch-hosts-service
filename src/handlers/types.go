package handlers

import (
	"bytes"
	"net/http"
	"sync"
)

type Channel struct {
	data     map[string]chan bytes.Buffer
	com      map[string]chan struct{}
	mux      *sync.Mutex
	chanSize uint
}

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

type Handler interface {
	HandleConsumer(request *http.Request, responseWriter http.ResponseWriter)
	HandleProducer(request *http.Request, responseWriter http.ResponseWriter)
}

var pubSubChannels = &Channel{data: make(map[string]chan bytes.Buffer), com: make(map[string]chan struct{}), mux: &sync.Mutex{}, chanSize: 100}
var defaultChannels = &Channel{data: make(map[string]chan bytes.Buffer), com: make(map[string]chan struct{}), mux: &sync.Mutex{}}

func NewHandlerStandard(urlPath string) Handler {
	data, _ := defaultChannels.getChannels(urlPath)
	return HandlerStandard{data: data}
}

func NewHandlerPubSub(urlPath string) Handler {
	data, com := pubSubChannels.getChannels(urlPath)
	return HandlerPubSub{data: data, com: com}
}
