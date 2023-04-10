package main

import (
	"io"
	"net/http"
)

func main() {

	r, err := DirectGet()
	if err != nil {
		println("error: ", err.Error())
		return
	}
	println("received response: ", string(r))
}

func DirectGet() ([]byte, error) {
	res, err := http.Get("https://httpbin.org/get")
	if err != nil {
		println("Failed to make request err: ", err)
		return nil, err
	}
	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		println("Failed to read body err: ", err)
		return nil, err
	}
	defer res.Body.Close()
	return resByte, nil
}
