/*
Copyright © 2025 PATRICK HERMANN patrick.hermann@sva.de
*/

package kubernetes

import (
	"fmt"
	"io"
	"net/http"
)

func FetchYAML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("FAILED TO FETCH YAML: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("FAILED TO FETCH YAML, STATUS CODE: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("FAILED TO READ YAML: %w", err)
	}

	// RETURN YAML CONTENT AS STRING
	return string(data), nil
}
