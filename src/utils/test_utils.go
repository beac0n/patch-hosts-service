package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Assert(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("actual value was %v, expected %v", actual, expected)
	}
}

func AssertRequest(expected interface{}, requestRecorder *httptest.ResponseRecorder, test *testing.T) bool {
	if status := requestRecorder.Code; status != http.StatusOK {
		test.Errorf("response has wrong status code: got %v expected %v ", status, http.StatusOK)
		return false
	}

	if expected == nil {
		return true
	}

	if actual := requestRecorder.Body.String(); actual != expected {
		test.Errorf("response has unexpected body: got %v expected %v", actual, expected)
		return false
	}

	return true
}

func SendRequest(reqHandler http.Handler, req *http.Request, reqRecChan chan *httptest.ResponseRecorder) {
	requestRecord := httptest.NewRecorder()

	reqHandler.ServeHTTP(requestRecord, req)

	reqRecChan <- requestRecord
}

func SendRequestSync(requestHandler http.Handler, postRequest *http.Request) *httptest.ResponseRecorder {
	// we need to wait a minimum amount of time, or else requests will be faster than writing to channel
	time.Sleep(time.Millisecond)

	requestRecorderPost := httptest.NewRecorder()
	requestHandler.ServeHTTP(requestRecorderPost, postRequest)
	return requestRecorderPost
}

func SendRequestWithCom(reqHandler http.Handler, req *http.Request, reqRecChan chan *httptest.ResponseRecorder, comChan chan struct{}) {
	requestRecord := httptest.NewRecorder()

	comChan <- struct{}{}

	reqHandler.ServeHTTP(requestRecord, req)

	reqRecChan <- requestRecord
}
