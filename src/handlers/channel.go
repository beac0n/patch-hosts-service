package handlers

import (
	"bytes"
	"sync"
)

type Channel struct {
	data *mapChanBuffer
	com  *mapChanStruct
}

type mapChanBuffer struct {
	m *sync.Map
}

func (m *mapChanBuffer) LoadOrStore(path string, channel chan bytes.Buffer) (chan bytes.Buffer, bool) {
	actual, loaded := m.m.LoadOrStore(path, channel)
	return actual.(chan bytes.Buffer), loaded
}

type mapChanStruct struct {
	m *sync.Map
}

func (m *mapChanStruct) LoadOrStore(path string, channel chan struct{}) (chan struct{}, bool) {
	actual, loaded := m.m.LoadOrStore(path, channel)
	return actual.(chan struct{}), loaded
}

func (channel *Channel) getChannels(path string) (chan bytes.Buffer, chan struct{}) {
	data, _ := channel.data.LoadOrStore(path, make(chan bytes.Buffer))
	com, _ := channel.com.LoadOrStore(path, make(chan struct{}))

	return data, com
}

func newChannel() *Channel {
	return &Channel{data: &mapChanBuffer{&sync.Map{}}, com: &mapChanStruct{&sync.Map{}}}
}
