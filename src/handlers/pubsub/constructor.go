package pubsub

type Handler struct {
	data           chan *[]byte
	com            chan struct{}
	maxReqSizeInMb int64
}

var pubSubChannel = newChannel()
var defaultChannel = newChannel()

func NewHandlerStandard(urlPath string, maxReqSizeInMb int64) Handler {
	data, _ := defaultChannel.getChannels(urlPath)
	return Handler{data: data, maxReqSizeInMb: maxReqSizeInMb}
}

func NewHandlerPubSub(urlPath string, maxReqSizeInMb int64) Handler {
	data, com := pubSubChannel.getChannels(urlPath)
	return Handler{data: data, com: com, maxReqSizeInMb: maxReqSizeInMb}
}
