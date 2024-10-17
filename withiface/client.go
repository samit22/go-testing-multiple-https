package withiface

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// Requester has Do method to make the request
type Requester interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPRequest to make http calls
type HTTPRequest struct {
	Client Requester
}

// GetHTTPBin gets data from http bin
func (h *HTTPRequest) GetHTTPBin(ctx context.Context) ([]byte, error) {
	url := "https://httpbin.org/get"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := h.Client.Do(req)
	if err != nil {
		return nil, err
	}
	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return resByte, fmt.Errorf("http status code %d", res.StatusCode)
	}
	return resByte, nil
}
