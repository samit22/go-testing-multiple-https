package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	_, err := Client("https://httpbin.org")
	if err != nil {
		println("Failed to make request err: ", err)
	}
}

func Client(baseURL string) ([]byte, error) {
	url := baseURL + "/get"
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: time.Second * 10}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}
	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resByte, nil
}
