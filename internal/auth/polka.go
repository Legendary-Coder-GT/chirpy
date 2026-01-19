package auth

import (
	"net/http"
	"strings"
	"errors"
)

func GetAPIKey(headers http.Header) (string, error) {
	if value, ok := headers["Authorization"]; ok {
		trimmed := strings.TrimSpace(value[0])
		words := strings.Split(trimmed, " ")
		return words[1], nil
	}
	return "", errors.New("No authorization header found")
}