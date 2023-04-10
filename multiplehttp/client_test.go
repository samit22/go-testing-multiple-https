package multiplehttp

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockReq struct {
	method  string
	url     string
	payload []byte
}

type mockRes struct {
	response   []byte
	statusCode int
	err        error
}

type mockClient struct {
	called int
	req    map[int]mockReq
	res    map[int]mockRes
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	m.called++

	r := mockReq{
		method: req.Method,
		url:    req.URL.String(),
	}
	if req.Body != nil {
		b, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		payload, err := io.ReadAll(b)
		if err != nil {
			return nil, err
		}
		r.payload = payload
	}
	m.req[m.called] = r

	mockResponse := m.res[m.called]
	respBody := io.NopCloser(bytes.NewReader(mockResponse.response))
	res := &http.Response{
		StatusCode: mockResponse.statusCode,
		Body:       respBody,
	}
	return res, mockResponse.err
}

func Test_GetHTTPBin(t *testing.T) {
	t.Log("successful response")
	{
		client := &mockClient{
			req: map[int]mockReq{},
			res: map[int]mockRes{
				1: {
					statusCode: 200,
					response:   []byte(`{"uuid": "3c95e984-b50c-471b-8f67-c2ace3809b06"}`),
				},
				2: {
					statusCode: 200,
					response:   []byte(`{"message": "all good from post"}`),
				},
				3: {
					statusCode: 200,
					response:   []byte(`{"message": "all good from get"}`),
				},
			},
		}

		req := HTTPRequest{
			Client: client,
		}

		_, err := req.MultipleHTTP()
		assert.NoError(t, err, "unexpected error %v", err)

		assert.Equal(t, 3, client.called, "expected to have called once.")

		// testing first http call
		assert.Equal(t, "GET", client.req[1].method)
		assert.Equal(t, "https://httpbin.org/uuid", client.req[1].url)

		// testing second http call
		assert.Equal(t, "POST", client.req[2].method)
		assert.Equal(t, "https://httpbin.org/post", client.req[2].url)
		assert.Equal(t, []byte(`{"id":"3c95e984-b50c-471b-8f67-c2ace3809b06"}`), client.req[2].payload)

		// testing third http call
		assert.Equal(t, "GET", client.req[3].method)
		assert.Equal(t, "https://httpbin.org/get", client.req[3].url)
	}
}
