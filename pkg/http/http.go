package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Get[T any](url string) (T, error) {
	var empty T
	resp, err := http.Get(url)
	if err != nil {
		return empty, fmt.Errorf("failed to send http get request: %v", err)
	}
	defer resp.Body.Close()

	var result T
	if dErr := json.NewDecoder(resp.Body).Decode(&result); dErr != nil {
		return empty, fmt.Errorf("failed to decoding response %v", dErr)
	}

	return result, nil
}
