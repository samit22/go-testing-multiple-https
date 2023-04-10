package httpclient

import (
	"io"
	"net/http"
	"time"
)

func Client() ([]byte, error) {
	url := "https://httpbin.org/get"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: time.Second * 10}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return resByte, nil
}
