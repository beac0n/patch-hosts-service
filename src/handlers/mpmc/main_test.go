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

func TestReqHandler_ServeHTTP_GET_POST(t *testing.T) {
	getReq := httptest.NewRequest("GET", "/foobar", nil)
	postReq := httptest.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(testData0)))

	reqRecChan := make(chan *httptest.ResponseRecorder)

	go utils.SendRequest(mpmcReqHandler, getReq, reqRecChan)

	utils.AssertRequest(nil, utils.SendRequestSync(mpmcReqHandler, postReq), t)
	utils.AssertRequest(testData0, <-reqRecChan, t)
}

func TestReqHandler_ServeHTTP_POST_GET(t *testing.T) {
	getReq := httptest.NewRequest("GET", "/foobar", nil)
	postReq := httptest.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(testData0)))

	reqRecChan := make(chan *httptest.ResponseRecorder)

	go utils.SendRequest(mpmcReqHandler, postReq, reqRecChan)

	utils.AssertRequest(testData0, utils.SendRequestSync(mpmcReqHandler, getReq), t)
	utils.AssertRequest(nil, <-reqRecChan, t)
}

func TestReqHandler_ServeHTTP_wrongMethod(t *testing.T) {
	postReq := httptest.NewRequest("PUT", "/foobar", bytes.NewBuffer([]byte("")))

	recorder := utils.SendRequestSync(mpmcReqHandler, postReq)
	utils.Assert(t, recorder.Code, http.StatusBadRequest)
}

func TestReqHandler_ServeHTTP_contentLengthError(t *testing.T) {
	postReq := httptest.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte("")))

	recorder := utils.SendRequestSync(mpmcReqHandler, postReq)
	utils.Assert(t, recorder.Code, http.StatusBadRequest)
}

func TestReqHandler_ServeHTTP_reqEntityTooLarge(t *testing.T) {
	mpmcReqHandler := NewReqHandler(1)
	postReq := httptest.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte("toolarge")))

	recorder := utils.SendRequestSync(mpmcReqHandler, postReq)
	utils.Assert(t, recorder.Code, http.StatusRequestEntityTooLarge)
}

func TestReqHandler_ServeHTTP_multi(t *testing.T) {
	getReq0 := httptest.NewRequest("GET", "/foobar", nil)
	postReq0 := httptest.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(testData0)))

	getReq1 := httptest.NewRequest("GET", "/barfoo", nil)
	postReq1 := httptest.NewRequest("POST", "/barfoo", bytes.NewBuffer([]byte(testData1)))

	reqRecChan0 := make(chan *httptest.ResponseRecorder)
	reqRecChan1 := make(chan *httptest.ResponseRecorder)

	go utils.SendRequest(mpmcReqHandler, getReq0, reqRecChan0)
	go utils.SendRequest(mpmcReqHandler, getReq1, reqRecChan1)

	utils.AssertRequest(nil, utils.SendRequestSync(mpmcReqHandler, postReq0), t)
	utils.AssertRequest(testData0, <-reqRecChan0, t)

	utils.AssertRequest(nil, utils.SendRequestSync(mpmcReqHandler, postReq1), t)
	utils.AssertRequest(testData1, <-reqRecChan1, t)
}
