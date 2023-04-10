package withiface

import (
	"bytes"
	"errors"
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
	req mockReq
	res mockRes
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {

	m.req = mockReq{
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
		m.req.payload = payload
	}

	r := io.NopCloser(bytes.NewReader(m.res.response))
	res := &http.Response{
		StatusCode: m.res.statusCode,
		Body:       r,
	}
	return res, m.res.err
}

func Test_GetHTTPBin(t *testing.T) {
	t.Log("successful response")
	{
		client := &mockClient{
			res: mockRes{
				statusCode: 200,
				response:   []byte(`{"message": "I am okay"}`)},
		}

		req := HTTPRequest{
			Client: client,
		}

		res, err := req.GetHTTPBin()
		assert.NoError(t, err, "unexpected error %v", err)
		assert.Equal(t, client.res.response, res)

		// test request
		assert.Equal(t, client.req.method, "GET")
		assert.Equal(t, client.req.url, "https://httpbin.org/get")

	}

	t.Log("mock client error")
	{
		client := &mockClient{
			res: mockRes{err: errors.New("something went wrong")},
		}

		req := HTTPRequest{
			Client: client,
		}

		_, err := req.GetHTTPBin()
		assert.NotNil(t, err, "should have got error")
	}

	t.Log("not okay http status code")
	{
		client := &mockClient{
			res: mockRes{statusCode: 400},
		}

		req := HTTPRequest{
			Client: client,
		}

		_, err := req.GetHTTPBin()
		assert.NotNil(t, err)
		assert.Equal(t, "http status code 400", err.Error())
	}
}
