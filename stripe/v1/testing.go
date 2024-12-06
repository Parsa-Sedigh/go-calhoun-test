package stripe

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient(t *testing.T) (*Client, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	c := &Client{
		baseURL: server.URL,
	}

	// returning func is a teardown func
	return c, mux, func() {
		server.Close()
	}
}
