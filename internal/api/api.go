package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func GET[T any](url string) (*T, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return call[T](req)
}

func call[T any](req *http.Request) (*T, error) {
	client := &http.Client{}
	userAgent := fmt.Sprintf("farstu/%s", "0.0.0")
	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	msg := fmt.Sprintf("%s %s %s %s", req.Proto, req.Method, res.Status, req.URL.String())
	slog.Debug(msg,
		"status", res.Status,
		"method", req.Method,
		"url", req.URL.String())

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%s - %s %s returned the following error:\n%s\n", res.Status, req.Method, req.URL.String(), string(bodyBytes))
	}

	var response T
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
