package pubsub

type Handler struct {
	data           chan *[]byte
	com            chan struct{}
	maxReqSizeInMb int64
}

var multiChannel = NewChannel()
var singlesChannel = NewChannel()

func NewHandlerSingle(urlPath string, maxReqSizeInMb int64) Handler {
	data, _ := singlesChannel.getChannels(urlPath)
	return Handler{data: data, maxReqSizeInMb: maxReqSizeInMb}
}

func NewHandlerMulti(urlPath string, maxReqSizeInMb int64) Handler {
	data, com := multiChannel.getChannels(urlPath)
	return Handler{data: data, com: com, maxReqSizeInMb: maxReqSizeInMb}
}
