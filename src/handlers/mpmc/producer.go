package mpmc

import (
	"io/ioutil"
	"net/http"
	"patch-hosts-service/src/constants"
	"patch-hosts-service/src/utils"
)

func (reqHandler *ReqHandler) produce(req *http.Request, resWriter http.ResponseWriter, dataChan chan *[]byte) {
	if req.ContentLength <= 0 {
		http.Error(resWriter, constants.NoContent, http.StatusBadRequest)
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
