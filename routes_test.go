package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode is a simple utility to check the response code
// of the response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestSignIn(t *testing.T) {
	s := NewServer()
	s.MountHandlers()

	req, _ := http.NewRequest("POST", "/signin", bytes.NewBufferString(`{
		"password": "password1",
		"username": "user1"
	}`))
	response := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestWelcome(t *testing.T) {
	s := NewServer()
	s.MountHandlers()

	sReq, _ := http.NewRequest("POST", "/signin", bytes.NewBufferString(`{
		"password": "password1",
		"username": "user1"
	}`))

	sResponse := executeRequest(sReq, s)

	req, _ := http.NewRequest("GET", "/welcome", nil)
	req.Header.Set("Auth-Token", sResponse.Body.String())
	response := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, response.Code)

	require.Equal(t, "welcome", response.Body.String())
}
