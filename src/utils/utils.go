package utils

import (
	"log"
	"net/http"
	"patch-hosts-service/src/constants"
	"runtime/debug"
	"strconv"
	"strings"
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

func NotGetOrPost(req *http.Request, resWriter http.ResponseWriter) bool {
	if (req.Method != http.MethodGet) && (req.Method != http.MethodPost) {
		http.Error(resWriter, constants.WrongHttpMethod, http.StatusBadRequest)
		return true
	}

	return false
}

func IsCorrectPath(request *http.Request, path string) bool {
	return strings.HasPrefix(request.URL.Path, path) && request.URL.Path != path
}
