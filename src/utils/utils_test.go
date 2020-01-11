package utils

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

var request = httptest.NewRequest("POST", "/test/path", bytes.NewBuffer([]byte("test_data")))

func TestLogErrorFalse(t *testing.T) {
	Assert(t, LogError(nil, request), false)
}

func TestLogErrorTrue(t *testing.T) {
	Assert(t, LogError(errors.New("t error"), request), true)
}

func TestHttpErrorRequestEntityTooLargeTrue(t *testing.T) {
	Assert(t, HttpErrorRequestEntityTooLarge(1, request, httptest.NewRecorder()), true)

}

func TestHttpErrorRequestEntityTooLargeFalse(t *testing.T) {
	Assert(t, HttpErrorRequestEntityTooLarge(10000, request, httptest.NewRecorder()), false)
}

func TestLoadAndStore(t *testing.T) {
	key := "test_key"
	m := &sync.Map{}
	chanCreator := func() interface{} { return make(chan struct{}) }

	channel := LoadAndStore(m, key, chanCreator)

	Assert(t, reflect.TypeOf(channel), reflect.TypeOf(make(chan struct{})))
}

func TestNotGetOrPostTrue(t *testing.T) {
	putRequest := httptest.NewRequest("PUT", "/t/path", nil)
	Assert(t, NotGetOrPost(putRequest, httptest.NewRecorder()), true)
}

func TestNotGetOrPostFalse(t *testing.T) {
	Assert(t, NotGetOrPost(request, httptest.NewRecorder()), false)
}

func TestIsCorrectPathTrue(t *testing.T) {
	Assert(t, IsCorrectPath(request,"/t"), true)
}

func TestIsCorrectPathFalsePrefix(t *testing.T) {
	Assert(t, IsCorrectPath(request,"/not-t"), false)
}

func TestIsCorrectPathFalse(t *testing.T) {
	Assert(t, IsCorrectPath(request,"/t/path"), false)
}
