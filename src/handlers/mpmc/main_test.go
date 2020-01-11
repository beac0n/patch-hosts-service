package mpmc

import (
	"../../utils"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mpmcReqHandler = NewReqHandler(10)
var testData0 = "test"
var testData1 = "test2"

func TestServeHttpSingle(t *testing.T) {
	getReq := httptest.NewRequest("GET", "/foobar", nil)
	postReq := httptest.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(testData0)))

	reqHandler := http.HandlerFunc(mpmcReqHandler.ServeHTTP)

	reqRecChan := make(chan *httptest.ResponseRecorder)

	go utils.SendRequest(reqHandler, getReq, reqRecChan)

	utils.AssertRequest("", utils.SendRequestSync(reqHandler, postReq), t)
	utils.AssertRequest(testData0, <-reqRecChan, t)
}

func TestServeHttpSingleParallel(t *testing.T) {
	getReq0 := httptest.NewRequest("GET", "/foobar", nil)
	postReq0 := httptest.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(testData0)))

	getReq1 := httptest.NewRequest("GET", "/barfoo", nil)
	postReq1 := httptest.NewRequest("POST", "/barfoo", bytes.NewBuffer([]byte(testData1)))

	reqHandler := http.HandlerFunc(mpmcReqHandler.ServeHTTP)

	reqRecChan0 := make(chan *httptest.ResponseRecorder)
	reqRecChan1 := make(chan *httptest.ResponseRecorder)

	go utils.SendRequest(reqHandler, getReq0, reqRecChan0)
	go utils.SendRequest(reqHandler, getReq1, reqRecChan1)

	utils.AssertRequest("", utils.SendRequestSync(reqHandler, postReq0), t)
	utils.AssertRequest(testData0, <-reqRecChan0, t)

	utils.AssertRequest("", utils.SendRequestSync(reqHandler, postReq1), t)
	utils.AssertRequest(testData1, <-reqRecChan1, t)
}
