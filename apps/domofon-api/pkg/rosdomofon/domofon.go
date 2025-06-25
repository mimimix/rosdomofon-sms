package rosdomofon

import (
	"fmt"
	"net/url"

	"domofon-api.gg/config"

	"github.com/imroc/req/v3"
)

type Domofon struct {
	accessToken  string
	refreshToken string
	baseURL      string
	client       *req.Client
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	UID          string `json:"uid"`
	Phone        string `json:"phone"`
	Company      string `json:"company"`
	JTI          string `json:"jti"`
}

type TemporaryKeyResponse struct {
	ActivationLink string `json:"activationLink"`
}

func NewDomofon(config *config.Config) *Domofon {
	return &Domofon{
		refreshToken: config.RefreshToken,
		baseURL:      "https://rdba.rosdomofon.com",
		client:       req.C().SetBaseURL("https://rdba.rosdomofon.com"),
	}
}

func (d *Domofon) refreshAccessToken() error {
	formData := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     "abonent",
		"refresh_token": d.refreshToken,
	}

	var tokenResp TokenResponse
	resp, err := d.client.R().
		SetFormData(formData).
		SetSuccessResult(&tokenResp).
		Post("/authserver-service/oauth/token")

	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	if !resp.IsSuccessState() {
		return fmt.Errorf("failed to refresh token, status: %s", resp.Status)
	}

	d.accessToken = "Bearer " + tokenResp.AccessToken
	d.refreshToken = tokenResp.RefreshToken
	return nil
}

func (d *Domofon) CreateTemporaryKey(KeyID int) (string, error) {
	var result TemporaryKeyResponse

	type RequestBody struct {
		ActivationsCount int `json:"activationsCount"`
		KeyID            int `json:"keyId"`
		WorkingPeriod    int `json:"workingPeriod"`
	}

	body := RequestBody{
		ActivationsCount: 1,
		KeyID:            KeyID,
		WorkingPeriod:    12,
	}

	if d.accessToken == "" {
		if err := d.refreshAccessToken(); err != nil {
			return "", fmt.Errorf("failed to refresh access token: %w", err)
		}
	}

	if d.accessToken == "" {
		if err := d.refreshAccessToken(); err != nil {
			return "", fmt.Errorf("failed to refresh access token: %w", err)
		}
	}

	sendRequest := func() (*req.Response, error) {
		return d.client.R().
			SetHeader("Authorization", d.accessToken).
			SetHeader("Content-Type", "application/json").
			SetBody(body).
			SetSuccessResult(&result).
			Post("/rdas-service/api/v1/temporary_keys")
	}

	// First attempt
	resp, err := sendRequest()
	if err != nil {
		return "", fmt.Errorf("failed to create temporary key: %w", err)
	}

	// If 401 Unauthorized, try to refresh token and retry
	if resp.StatusCode == 401 {
		if err := d.refreshAccessToken(); err != nil {
			return "", fmt.Errorf("failed to refresh access token: %w", err)
		}

		// Retry after token refresh
		resp, err = sendRequest()
		if err != nil {
			return "", fmt.Errorf("failed to create temporary key after token refresh: %w", err)
		}
	}

	if !resp.IsSuccessState() {
		return "", fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return result.ActivationLink, nil
}

// ActivateKey activates a temporary key using the provided activation link
// The activationLink should be in format: https://my.rosdomofon.com/temporary-keys/activate?token=TOKEN_VALUE
// Returns nil on success (HTTP 204), or an error if activation fails
func (d *Domofon) ActivateKey(activationLink string) error {
	// Parse the URL to extract the token
	parsedURL, err := url.Parse(activationLink)
	if err != nil {
		return fmt.Errorf("invalid activation link: %w", err)
	}

	// Extract the token from query parameters
	token := parsedURL.Query().Get("token")
	if token == "" {
		return fmt.Errorf("no token found in activation link")
	}

	// Build the activation URL
	activationURL := fmt.Sprintf("/rdas-service/api/v1/temporary_keys/%s/activate", token)

	// Send the activation request
	resp, err := d.client.R().
		Post(activationURL)

	if err != nil {
		return fmt.Errorf("failed to activate key: %w", err)
	}

	// Check if we got a 204 No Content response
	if resp.StatusCode != 204 {
		return fmt.Errorf("activation failed with status: %s", resp.Status)
	}

	return nil
}
