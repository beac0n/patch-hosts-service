package mpmc

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mpmcReqHandler = NewReqHandler(10)
var testData0 = "test"
var testData1 = "test2"

func TestServeHttpSingle(test *testing.T) {
	getReq, _ := http.NewRequest("GET", "/foobar", nil)
	postReq, _ := http.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(testData0)))

	reqHandler := http.HandlerFunc(mpmcReqHandler.ServeHTTP)

	reqRecChan := make(chan *httptest.ResponseRecorder)

	go sendReq(reqHandler, getReq, reqRecChan)

	assertReq("", sendReqSync(reqHandler, postReq), test)
	assertReq(testData0, <-reqRecChan, test)
}

func TestServeHttpSingleParallel(test *testing.T) {
	getReq0, _ := http.NewRequest("GET", "/foobar", nil)
	postReq0, _ := http.NewRequest("POST", "/foobar", bytes.NewBuffer([]byte(testData0)))

	getReq1, _ := http.NewRequest("GET", "/barfoo", nil)
	postReq1, _ := http.NewRequest("POST", "/barfoo", bytes.NewBuffer([]byte(testData1)))

	reqHandler := http.HandlerFunc(mpmcReqHandler.ServeHTTP)

	reqRecChan0 := make(chan *httptest.ResponseRecorder)
	reqRecChan1 := make(chan *httptest.ResponseRecorder)

	go sendReq(reqHandler, getReq0, reqRecChan0)
	go sendReq(reqHandler, getReq1, reqRecChan1)

	assertReq("", sendReqSync(reqHandler, postReq0), test)
	assertReq(testData0, <-reqRecChan0, test)

	assertReq("", sendReqSync(reqHandler, postReq1), test)
	assertReq(testData1, <-reqRecChan1, test)
}

func sendReqSync(requestHandler http.HandlerFunc, postRequest *http.Request) *httptest.ResponseRecorder {
	requestRecorderPost := httptest.NewRecorder()
	requestHandler.ServeHTTP(requestRecorderPost, postRequest)
	return requestRecorderPost
}

func assertReq(expected string, requestRecorder *httptest.ResponseRecorder, test *testing.T) {
	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("response has wrong status code: got %v want %v ", status, http.StatusOK)
	}

	if expected == "" {
		return
	}

	if actual := requestRecorder.Body.String(); actual != expected {
		test.Errorf("response has unexpected body: got %v want %v", actual, expected)
	}
}

func sendReq(reqHandler http.HandlerFunc, req *http.Request, reqRecChan chan *httptest.ResponseRecorder) {
	requestRecord := httptest.NewRecorder()

	reqHandler.ServeHTTP(requestRecord, req)

	reqRecChan <- requestRecord
}
