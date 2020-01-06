package mpmc

import (
	"../../utils"
	"io/ioutil"
	"net/http"
)

func (reqHandler *ReqHandler) produce(req *http.Request, resWriter http.ResponseWriter, dataChan chan *[]byte) {
	if req.ContentLength <= 0 {
		http.Error(resWriter, "no content", http.StatusBadRequest)
		return
	}

	if utils.HttpErrorRequestEntityTooLarge(reqHandler.maxReqSize, req, resWriter) {
		return
	}

	bodyBytes, err := ioutil.ReadAll(req.Body)

	if utils.LogError(err, req) {
		return
	}

	select {
	case dataChan <- &bodyBytes:
	case <-req.Context().Done():
	}

}
