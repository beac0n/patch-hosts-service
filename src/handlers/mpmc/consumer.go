package mpmc

import (
	"../../utils"
	"net/http"
	"strconv"
)

func (reqHandler *RequestHandler) consume(req *http.Request, resWriter http.ResponseWriter, dataChan chan *[]byte) {
	select {
	case bytes := <-dataChan:
		resWriter.Header().Set("Content-Length", strconv.Itoa(len(*bytes)))
		_, err := resWriter.Write(*bytes)
		utils.LogError(err, req)

	case <-req.Context().Done():
	}
}
