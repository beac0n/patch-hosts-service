package pubsub

import "sync"

type dataChannelMap struct {
	syncMap *sync.Map
}

func (mapChanBuffer *dataChannelMap) LoadOrStore(path string, channel chan dataChanel) chan dataChanel {
	actual, _ := mapChanBuffer.syncMap.LoadOrStore(path, channel)
	return actual.(chan dataChanel)
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

func (channel *channelMapWrap) getDataChannel(path string) chan dataChanel {
	return channel.data.LoadOrStore(path, make(chan dataChanel))
}

func (channel *channelMapWrap) getComChannel(path string) chan struct{} {
	return channel.com.LoadOrStore(path, make(chan struct{}))
}

func newChannelMapWrap() *channelMapWrap {
	return &channelMapWrap{
		data: &dataChannelMap{&sync.Map{}},
		com:  &comChannelMap{&sync.Map{}},
	}
}
