package pubsub

import (
	"../../utils"
	"net/http"
)

func (requestHandler *RequestHandler) consume(request *http.Request, responseWriter http.ResponseWriter, dataChannel chan *[]byte, comChannel chan struct{}, persist bool) {
	comChannel <- struct{}{}

	for {
		select {
		case bytes := <-dataChannel:
			if persist {
				comChannel <- struct{}{}
			}

			_, err := responseWriter.Write(*bytes)
			responseWriter.(http.Flusher).Flush()

			if utils.LogError(err, request) || !persist {
				return
			}

		case <-request.Context().Done():
			<-comChannel
			return
		}
	}
}
