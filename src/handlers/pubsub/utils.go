package pubsub

import "net/http"

func (channelWrap channelWrap) sendDataToConsumers(consumersCount uint64, bodyBytes []byte, bodyLen int64, request *http.Request) {
	for ; consumersCount > 0; consumersCount-- {
		select {
		case channelWrap.data <- dataChanel{&bodyBytes, bodyLen}:
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
ComLoop:
	for {
		select {
		case <-channelWrap.com:
			consumersCount++
		default:
			break ComLoop
		}
	}

	return consumersCount
}
