package plainhttp

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	_, err := DirectGet("https://httpbin.org")
	if err != nil {
		fmt.Printf("Failed to make request err: %v", err)
	}
}

func DirectGet(baseURL string) ([]byte, error) {
	url := baseURL + "/get"
	res, err := http.Get(url)
	if err != nil {
		println("Failed to make request err: ", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		println("Failed to make request err: ", err)
		return nil, fmt.Errorf("bad response status %s", res.Status)
	}
	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		println("Failed to read body err: ", err)
		return nil, err
	}
	return resByte, nil
}
