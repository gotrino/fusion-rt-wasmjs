package auth

import "net/http"

type Authenticator interface {
	Authenticate(request http.Request)
}

// bearerTokenResponse see also https://datatracker.ietf.org/doc/html/rfc6750#section-4
type bearerTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
