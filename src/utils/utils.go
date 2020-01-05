package utils

import (
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"sync"
)

func LogError(err error, request *http.Request) bool {
	if err == nil {
		return false
	}

	log.Println("ERROR", request.Method, request.URL.Path, err)
	debug.PrintStack()

	return true

}

func HttpErrorRequestEntityTooLarge(maxReqSize int64, request *http.Request, responseWriter http.ResponseWriter) bool {
	if request.ContentLength <= maxReqSize {
		return false
	}

	maxReqSizeInByteStr := strconv.FormatInt(maxReqSize, 10)
	reqContentLenStr := strconv.FormatInt(request.ContentLength, 10)
	errorMsg := "max. request size is " + maxReqSizeInByteStr + ", got " + reqContentLenStr
	http.Error(responseWriter, errorMsg, http.StatusRequestEntityTooLarge)

	return true
}

func LoadAndStore(syncMap *sync.Map, key string, channelCreator func() interface{}) interface{} {
	if _, dataChannelOk := syncMap.Load(key); !dataChannelOk {
		syncMap.Store(key, channelCreator())
	}

	dataChannelI, _ := syncMap.Load(key)
	return dataChannelI
}