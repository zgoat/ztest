package test

import (
	"net/http/httptest"
	"testing"
)

func TestCode(t *testing.T) {
	Code(t, &httptest.ResponseRecorder{Code: 200}, 200)

	// TODO: how to test that t.Fatalf() was called?
	//TestCode(t, &httptest.ResponseRecorder{Code: 200}, 201)
}
