package main

import (
	"net/http"
	"net/http/httptest"
	"patch-hosts-service/src/utils"
	"testing"
)

type mockReqHandler struct {
	wasCalled bool
}

func (reqHandler *mockReqHandler) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	reqHandler.wasCalled = true
}

var resWriter = httptest.NewRecorder()

func TestReqHandler_ServeHTTPPubSub(t *testing.T) {
	req := httptest.NewRequest("GET", "/pubsub/foobar", nil)

	mockHandler := &mockReqHandler{}
	reqHandler := &ReqHandler{pubSubReqHandler: mockHandler}

	reqHandler.ServeHTTP(resWriter, req)

	utils.Assert(t, mockHandler.wasCalled, true)
}

func TestReqHandler_ServeHTTPQueue(t *testing.T) {
	req := httptest.NewRequest("GET", "/queue/foobar", nil)

	mockHandler := &mockReqHandler{}
	reqHandler := &ReqHandler{mpmcReqHandler: mockHandler}

	reqHandler.ServeHTTP(resWriter, req)

	utils.Assert(t, mockHandler.wasCalled, true)
}

func TestReqHandler_ServeHTTPReq(t *testing.T) {
	req := httptest.NewRequest("GET", "/req/foobar", nil)

	mockHandler := &mockReqHandler{}
	reqHandler := &ReqHandler{reqResReqHandler: mockHandler}

	reqHandler.ServeHTTP(resWriter, req)

	utils.Assert(t, mockHandler.wasCalled, true)
}

func TestReqHandler_ServeHTTPRes(t *testing.T) {
	req := httptest.NewRequest("GET", "/res/foobar", nil)

	mockHandler := &mockReqHandler{}
	reqHandler := &ReqHandler{reqResReqHandler: mockHandler}

	reqHandler.ServeHTTP(resWriter, req)

	utils.Assert(t, mockHandler.wasCalled, true)
}

func TestReqHandler_ServeHTTP404(t *testing.T) {
	req := httptest.NewRequest("GET", "/none/foobar", nil)
	resWriter := httptest.NewRecorder()

	(&ReqHandler{}).ServeHTTP(resWriter, req)

	utils.Assert(t, resWriter.Code, 404)
}
