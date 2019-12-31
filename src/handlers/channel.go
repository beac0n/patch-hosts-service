package handlers

import (
	"bytes"
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

func newChannel() *Channel {
	return &Channel{data: make(map[string]chan bytes.Buffer), com: make(map[string]chan struct{}), mux: &sync.Mutex{}}
}
