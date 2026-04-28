package main

import (
	"bytes"

	// "crypto/des"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping?user=jack", nil)

	w := httptest.NewRecorder()

	ping(w, req)

	desiredcode := http.StatusOK

	if desiredcode != w.Code {
		t.Errorf("Bad response expected: %v but received: %v", desiredcode, w.Code)
	}
	expectedMessage := []byte("Hello jack!\n")

	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("Bad response body expected: %q but received: %v", expectedMessage, w.Body.String())
	}
}
