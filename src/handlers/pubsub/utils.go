package pubsub

import (
	"net/http"
	"strconv"
)

func (channelWrap channelWrap) sendDataToConsumers(consumersCount uint64, bodyBytes *[]byte, request *http.Request) {
	for ; consumersCount > 0; consumersCount-- {
		select {
		case channelWrap.data <- bodyBytes:
		case <-request.Context().Done():
			return
		}
	}
}

func (channelWrap channelWrap) getConsumerCount() uint64 {
	if channelWrap.com == nil {
		return 1
	}

	consumersCount := uint64(0)

	for {
		select {
		case <-channelWrap.com:
			consumersCount++
		default:
			return consumersCount
		}
	}
}

func (channelWrap channelWrap) httpErrorEntityTooLarge(request *http.Request, responseWriter http.ResponseWriter) {
	maxReqSizeInByteStr := strconv.FormatInt(channelWrap.maxReqSize, 10)
	reqContentLenStr := strconv.FormatInt(request.ContentLength, 10)
	errorMsg := "max. request size is " + maxReqSizeInByteStr + ", got " + reqContentLenStr
	http.Error(responseWriter, errorMsg, http.StatusRequestEntityTooLarge)
}
