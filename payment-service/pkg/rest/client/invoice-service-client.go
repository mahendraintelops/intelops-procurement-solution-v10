package client

import (
	"io"
	"net/http"
)

func ExecuteNodeC2() ([]byte, error) {
	response, err := http.Get("invoice-service:3500/ping")

	if err != nil {
		return nil, err
	}

	return io.ReadAll(response.Body)
}
