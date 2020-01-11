package mpmc

import (
	"../../constants"
	"../../utils"
	"net/http"
	"strconv"
)

func (reqHandler *ReqHandler) consume(req *http.Request, resWriter http.ResponseWriter, dataChan chan *[]byte) {
	select {
	case bytes := <-dataChan:
		resWriter.Header().Set(constants.ContentLength, strconv.Itoa(len(*bytes)))
		_, err := resWriter.Write(*bytes)
		resWriter.(http.Flusher).Flush()
		utils.LogError(err, req)

	case <-req.Context().Done():
	}
}
