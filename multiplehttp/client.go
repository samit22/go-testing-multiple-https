package multiplehttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Requester has Do method to make the request
type Requester interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPRequest to make http calls
type HTTPRequest struct {
	Client Requester
}

func WithIface() {
	client := &http.Client{Timeout: time.Second * 10}
	r := HTTPRequest{
		Client: client,
	}
	r.MultipleHTTP()
}

// GetHTTPBin gets data from http bin
func (h *HTTPRequest) MultipleHTTP() ([]byte, error) {
	id, err := h.getUUID()
	if err != nil {
		return nil, err
	}
	err = h.postUUID(id)
	if err != nil {
		return nil, err
	}
	return h.normalGet()
}

func (h *HTTPRequest) getUUID() (string, error) {
	url := "https://httpbin.org/uuid"
	resp, err := h.makeHTTPCall("GET", url, nil)
	if err != nil {
		return "", err
	}
	type uuid struct {
		UUID string `json:"uuid"`
	}
	var u uuid
	err = json.Unmarshal(resp, &u)
	if err != nil {
		return "", err
	}
	return u.UUID, nil
}
func (h *HTTPRequest) postUUID(uuid string) error {
	url := "https://httpbin.org/post"
	body := map[string]string{"id": uuid}
	bt, _ := json.Marshal(body)
	_, err := h.makeHTTPCall("POST", url, bt)
	return err
}

func (h *HTTPRequest) normalGet() ([]byte, error) {
	url := "https://httpbin.org/get"
	return h.makeHTTPCall("GET", url, nil)
}

func (h *HTTPRequest) makeHTTPCall(method, url string, payload []byte) ([]byte, error) {
	client := h.Client
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("http status code %d", res.StatusCode)
	}
	return body, nil
}
