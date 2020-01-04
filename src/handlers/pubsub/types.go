package pubsub

import "sync"

type dataChannelMap struct {
	syncMap *sync.Map
}

func (mapChanBuffer *dataChannelMap) LoadOrStore(path string, channel chan *[]byte) chan *[]byte {
	actual, _ := mapChanBuffer.syncMap.LoadOrStore(path, channel)
	return actual.(chan *[]byte)
}

type comChannelMap struct {
	syncMap *sync.Map
}

func (mapChanStruct *comChannelMap) LoadOrStore(path string, channel chan struct{}) chan struct{} {
	actual, _ := mapChanStruct.syncMap.LoadOrStore(path, channel)
	return actual.(chan struct{})
}

type channelMapWrap struct {
	data *dataChannelMap
	com  *comChannelMap
}

func (channel *channelMapWrap) getDataChannel(path string, length uint) chan *[]byte {
	return channel.data.LoadOrStore(path, make(chan *[]byte, length))
}

func (channel *channelMapWrap) getComChannel(path string, length uint) chan struct{} {
	return channel.com.LoadOrStore(path, make(chan struct{}, length))
}

func newChannelMapWrap() *channelMapWrap {
	return &channelMapWrap{
		data: &dataChannelMap{&sync.Map{}},
		com:  &comChannelMap{&sync.Map{}},
	}
}
