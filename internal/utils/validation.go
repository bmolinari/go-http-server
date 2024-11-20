package utils

import (
	"fmt"
	"net/url"
)

func ValidateQueryParams(queryParams url.Values, requiredKeys []string) (map[string]string, error) {
	validatedParams := make(map[string]string)
	for _, key := range requiredKeys {
		value := queryParams.Get(key)
		if value == "" {
			return nil, fmt.Errorf("Missing required query parameter: '%s'", key)
		}
		validatedParams[key] = value
	}
	return validatedParams, nil
}

func ValidateContentType(headers map[string]string, expectedContentType string) bool {
	contentType, exists := headers["Content-Type"]
	if !exists {
		return false
	}
	return contentType == expectedContentType
}

func ValidateContentLength(headers map[string]string) (int, bool) {
	cl, exists := headers["Content-Length"]
	if !exists {
		return 0, false
	}

	var contentLength int
	_, err := fmt.Sscanf(cl, "%d", &contentLength)
	if err != nil || contentLength < 0 {
		return 0, false
	}

	return contentLength, true
}
