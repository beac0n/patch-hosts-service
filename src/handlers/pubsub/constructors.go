package pubsub

var multiChannel = newChannel()
var singlesChannel = newChannel()

type handler struct {
	data           chan *[]byte
	com            chan struct{}
	maxReqSizeInMb int64
}

func newHandlerSingle(urlPath string, maxReqSizeInMb int64) handler {
	data, _ := singlesChannel.getChannels(urlPath)
	return handler{data: data, maxReqSizeInMb: maxReqSizeInMb}
}

func newHandlerMulti(urlPath string, maxReqSizeInMb int64) handler {
	data, com := multiChannel.getChannels(urlPath)
	return handler{data: data, com: com, maxReqSizeInMb: maxReqSizeInMb}
}
