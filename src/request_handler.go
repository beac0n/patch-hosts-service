package main

import (
	"./constants"
	"./utils"
	"net/http"
)

type ReqHandler struct {
	pubSubReqHandler http.Handler
	mpmcReqHandler   http.Handler
	reqResReqHandler http.Handler
}

func (reqHandler *ReqHandler) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	if utils.IsCorrectPath(req, "/pubsub") {
		reqHandler.pubSubReqHandler.ServeHTTP(resWriter, req)
	} else if utils.IsCorrectPath(req, "/queue") {
		reqHandler.mpmcReqHandler.ServeHTTP(resWriter, req)
	} else if utils.IsCorrectPath(req, constants.Res) || utils.IsCorrectPath(req, constants.Req) {
		reqHandler.reqResReqHandler.ServeHTTP(resWriter, req)
	} else {
		http.Error(resWriter, "", http.StatusNotFound)
	}
}
