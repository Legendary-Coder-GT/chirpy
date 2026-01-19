package auth

import (
	"net/http"
	"errors"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	if value, ok := headers["Authorization"]; ok {
		words := strings.Split(value[0], " ")
		if len(words) != 2 || len(words[1]) == 0{
			return "", errors.New("Authorization header incorrectly formatted")
		}
		return words[1], nil
	}
	return "", errors.New("No authorization header found")
}