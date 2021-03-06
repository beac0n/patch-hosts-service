package pubsub

import (
	"net/http"
	"patch-hosts-service/src/utils"
)

func (reqHandler *ReqHandler) consume(req *http.Request, resWriter http.ResponseWriter, dataChan chan *[]byte, comChan chan struct{}, persist bool) {
	comChan <- struct{}{}

	for {
		select {
		case bytes := <-dataChan:
			if persist {
				comChan <- struct{}{}
			}

			_, err := resWriter.Write(*bytes)
			resWriter.(http.Flusher).Flush()

			if utils.LogError(err, req) || !persist {
				return
			}

		case <-req.Context().Done():
			<-comChan
			return
		}
	}
}
