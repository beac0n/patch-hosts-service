package handlers

import "bytes"

type Handler struct {
	data chan bytes.Buffer
	com  chan struct{}
}

var pubSubChannel = newChannel()
var defaultChannel = newChannel()

func NewHandlerStandard(urlPath string) Handler {
	data, _ := defaultChannel.getChannels(urlPath)
	return Handler{data: data}
}

func NewHandlerPubSub(urlPath string) Handler {
	data, com := pubSubChannel.getChannels(urlPath)
	return Handler{data: data, com: com}
}
