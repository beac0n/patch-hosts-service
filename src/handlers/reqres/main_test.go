package reqres

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"patch-hosts-service/src/utils"
	"testing"
)

var reqResReqHandler = NewReqHandler(10)
var testData = "test"

func TestReqHandler_ServeHTTP_GetGet(t *testing.T) {
	getReq0 := httptest.NewRequest("GET", "/req/foobar", bytes.NewBuffer([]byte(testData)))
	getReq1 := httptest.NewRequest("GET", "/res/foobar", bytes.NewBuffer([]byte(testData)))

	reqRecChan := make(chan *httptest.ResponseRecorder)

	go utils.SendRequest(reqResReqHandler, getReq0, reqRecChan)
	go utils.SendRequest(reqResReqHandler, getReq1, reqRecChan)

	utils.AssertRequest(testData, <-reqRecChan, t)
	utils.AssertRequest(testData, <-reqRecChan, t)
}

func TestReqHandler_ServeHTTP_PostPost(t *testing.T) {
	postReq0 := httptest.NewRequest("POST", "/req/foobar", bytes.NewBuffer([]byte(testData)))
	postReq1 := httptest.NewRequest("POST", "/res/foobar", bytes.NewBuffer([]byte(testData)))

	reqRecChan := make(chan *httptest.ResponseRecorder)

	go utils.SendRequest(reqResReqHandler, postReq0, reqRecChan)
	go utils.SendRequest(reqResReqHandler, postReq1, reqRecChan)

	utils.AssertRequest(testData, <-reqRecChan, t)
	utils.AssertRequest(testData, <-reqRecChan, t)
}

func TestReqHandler_ServeHTTP_PostPostWithExtraHeader(t *testing.T) {
	postReq0 := httptest.NewRequest("POST", "/req/foobar", bytes.NewBuffer([]byte(testData)))
	postReq0.Header.Set("foobar", "barfoo")

	postReq1 := httptest.NewRequest("POST", "/res/foobar", bytes.NewBuffer([]byte(testData)))
	postReq0.Header.Set("barfoo", "foobar")

	reqRecChan := make(chan *httptest.ResponseRecorder)

	go utils.SendRequest(reqResReqHandler, postReq0, reqRecChan)
	go utils.SendRequest(reqResReqHandler, postReq1, reqRecChan)

	recorderReq0 := <-reqRecChan
	utils.Assert(t, recorderReq0.Header().Get("X-Phs-0-Foobar"), "barfoo")
	utils.AssertRequest(testData, recorderReq0, t)

	recorderReq1 := <-reqRecChan
	utils.Assert(t, recorderReq0.Header().Get("X-Phs-0-Barfoo"), "foobar")
	utils.AssertRequest(testData, recorderReq1, t)
}

func TestReqHandler_ServeHTTP_reqEntityTooLarge(t *testing.T) {
	reqResReqHandler := NewReqHandler(1)
	postReq := httptest.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte("toolarge")))

	recorder := utils.SendRequestSync(reqResReqHandler, postReq)
	utils.Assert(t, recorder.Code, http.StatusRequestEntityTooLarge)
}

func TestReqHandler_ServeHTTP_unsupportedHttpMethod(t *testing.T) {
	postReq := httptest.NewRequest("PING", "/foobar", bytes.NewBuffer([]byte("")))

	recorder := utils.SendRequestSync(reqResReqHandler, postReq)
	utils.Assert(t, recorder.Code, http.StatusNotImplemented)
}
