package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type responseBody struct {
	Username string `json:"username"`
}

func VerifyCogToken(token string) (username string, err error) {
	if token == "" {
		return "", fmt.Errorf("token is required")
	}

	resp, err := http.PostForm("https://r8.im/cog/v1/verify-token", url.Values{
		"token": []string{token},
	})
	if err != nil {
		return "", fmt.Errorf("failed to verify token: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("user does not exist")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to verify token, got status %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body := &responseBody{}
	if err := json.NewDecoder(resp.Body).Decode(body); err != nil {
		return "", err
	}

	return body.Username, nil
}
