package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockServer struct {
	// return response
	respStatusCode int
	respBody       []byte

	// capture request
	reqPath   string
	reqMethod string
	reqHeader http.Header
}

func (m *mockServer) createServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.reqMethod = r.Method
		m.reqPath = r.URL.Path
		m.reqHeader = r.Header

		w.WriteHeader(m.respStatusCode)
		w.Write(m.respBody)
	}))
}

func TestDirectGet(t *testing.T) {
	t.Run("no error on request", func(t *testing.T) {
		mock := mockServer{
			respStatusCode: http.StatusOK,
			respBody:       []byte(`{"success": true}`),
		}
		mockServer := mock.createServer()
		defer func() { mockServer.Close() }()

		_, err := Client(mockServer.URL)
		if err != nil {
			t.Fatalf("Failed to make request err: %v", err)
		}

		if mock.reqMethod != "GET" {
			t.Errorf("Expected method GET, got %s", mock.reqMethod)
		}

		if mock.reqPath != "/get" {
			t.Errorf("Expected path /, got %s", mock.reqPath)
		}
	})

	t.Run("error on response status code", func(t *testing.T) {
		mock := mockServer{
			respStatusCode: http.StatusBadRequest,
			respBody:       []byte(`{"success": false}`),
		}
		mockServer := mock.createServer()
		defer func() { mockServer.Close() }()

		_, err := Client(mockServer.URL)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("error on making request", func(t *testing.T) {
		mock := mockServer{}
		mockServer := mock.createServer()
		defer func() { mockServer.Close() }()

		_, err := Client(mockServer.URL)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
