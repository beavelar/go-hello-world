package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHttpReq(t *testing.T) {
	request, _ := http.NewRequest("GET", "/http", nil)
	response := httptest.NewRecorder()

	httpReq(response, request)

	got := response.Body.String()
	want := "simple server http response"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGetRedirectedReq(t *testing.T) {
	request, _ := http.NewRequest("GET", "/redirected", nil)
	response := httptest.NewRecorder()

	redirectedReq(response, request)

	got := response.Body.String()
	want := "simple server redirected response"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGetRootReq(t *testing.T) {
	HttpServerConfig = &Config{Host: "localhost", Port: 80, Protocol: "http"}

	t.Run("root request with no redirect", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		rootReq(response, request)

		got := response.Body.String()
		want := "simple server root response"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("root request with redirect", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/test", nil)
		response := httptest.NewRecorder()

		rootReq(response, request)

		got := response.Code
		want := http.StatusPermanentRedirect

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
