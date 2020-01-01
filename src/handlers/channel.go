package handlers

import "sync"

type channel struct {
	data *mapChanBuffer
	com  *mapChanStruct
}

type mapChanBuffer struct {
	syncMap *sync.Map
}

func (mapChanBuffer *mapChanBuffer) LoadOrStore(path string, channel chan *[]byte) (chan *[]byte, bool) {
	actual, loaded := mapChanBuffer.syncMap.LoadOrStore(path, channel)
	return actual.(chan *[]byte), loaded
}

type mapChanStruct struct {
	syncMap *sync.Map
}

func (mapChanStruct *mapChanStruct) LoadOrStore(path string, channel chan struct{}) (chan struct{}, bool) {
	actual, loaded := mapChanStruct.syncMap.LoadOrStore(path, channel)
	return actual.(chan struct{}), loaded
}

func (channel *channel) getChannels(path string) (chan *[]byte, chan struct{}) {
	data, _ := channel.data.LoadOrStore(path, make(chan *[]byte))
	com, _ := channel.com.LoadOrStore(path, make(chan struct{}))

	return data, com
}

func newChannel() *channel {
	return &channel{data: &mapChanBuffer{&sync.Map{}}, com: &mapChanStruct{&sync.Map{}}}
}
