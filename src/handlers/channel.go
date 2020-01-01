package handlers

import "sync"

type Channel struct {
	data *mapChanBuffer
	com  *mapChanStruct
}

type mapChanBuffer struct {
	m *sync.Map
}

func (m *mapChanBuffer) LoadOrStore(path string, channel chan *[]byte) (chan *[]byte, bool) {
	actual, loaded := m.m.LoadOrStore(path, channel)
	return actual.(chan *[]byte), loaded
}

type mapChanStruct struct {
	m *sync.Map
}

func (m *mapChanStruct) LoadOrStore(path string, channel chan struct{}) (chan struct{}, bool) {
	actual, loaded := m.m.LoadOrStore(path, channel)
	return actual.(chan struct{}), loaded
}

func (channel *Channel) getChannels(path string) (chan *[]byte, chan struct{}) {
	data, _ := channel.data.LoadOrStore(path, make(chan *[]byte))
	com, _ := channel.com.LoadOrStore(path, make(chan struct{}))

	return data, com
}

func newChannel() *Channel {
	return &Channel{data: &mapChanBuffer{&sync.Map{}}, com: &mapChanStruct{&sync.Map{}}}
}
