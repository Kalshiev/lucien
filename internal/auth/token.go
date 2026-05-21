package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	if headers.Get("Authorization") == "" {
		return "", fmt.Errorf("No Authorization Token in Header")
	}

	token := strings.TrimPrefix(headers.Get("Authorization"), "Bearer ")

	return token, nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)

	return hex.EncodeToString(key)
}

func GetAPIKey(headers http.Header) (string, error) {
	if headers.Get("Authorization") == "" {
		return "", fmt.Errorf("No Authorization Token in Header")
	}

	key := strings.TrimPrefix(headers.Get("Authorization"), "ApiKey ")

	return key, nil
}
