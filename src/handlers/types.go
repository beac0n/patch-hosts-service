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
}

func (channel *Channel) getChannels(path string) (chan bytes.Buffer, chan struct{}) {
	channel.mux.Lock()
	defer channel.mux.Unlock()

	_, dataOk := channel.data[path]
	if !dataOk {
		channel.data[path] = make(chan bytes.Buffer)
	}

	_, comOk := channel.com[path]
	if !comOk {
		channel.com[path] = make(chan struct{})
	}

	return channel.data[path], channel.com[path]
}

type Handler interface {
	HandleConsumer(request *http.Request, responseWriter http.ResponseWriter)
	HandleProducer(request *http.Request, responseWriter http.ResponseWriter)
}

func newChannel() *Channel {
	return &Channel{data: make(map[string]chan bytes.Buffer), com: make(map[string]chan struct{}), mux: &sync.Mutex{}}
}

var pubSubChannel = newChannel()
var defaultChannel = newChannel()

func NewHandlerStandard(urlPath string) Handler {
	data, _ := defaultChannel.getChannels(urlPath)
	return HandlerPubSub{data: data}
}

func NewHandlerPubSub(urlPath string) Handler {
	data, com := pubSubChannel.getChannels(urlPath)
	return HandlerPubSub{data: data, com: com}
}
