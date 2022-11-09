package auth

import (
	"fmt"
)

// AuthMethod is a generic container producing token
type AuthMethod interface {
	Token() string
}

// Permit is the result of auth process. It contain secrets to be used in communication.
type Permit struct {
	token string
}

// Claim is a set of possible authorisation requisits. The auth method
// picks upsome of the fields. They are used to log in into the Azure AD.
type Claim struct {
	ClientID            string   `json:"client_id,omitempty"`
	Authority           string   `json:"authority,omitempty"`
	Scopes              []string `json:"scopes,omitempty"`
	Username            string   `json:"username,omitempty"`
	Password            string   `json:"password,omitempty"`
	RedirectURI         string   `json:"redirect_uri,omitempty"`
	CodeChallenge       string   `json:"code_challenge,omitempty"`
	CodeChallengeMethod string   `json:"code_challenge_method,omitempty"`
	State               string   `json:"state,omitempty"`
	ClientSecret        string   `json:"client_secret,omitempty"`
	Thumbprint          string   `json:"thumbprint,omitempty"`
	PemData             string   `json:"pem_data,omitempty"`
}

// NewAuthMethod is a factory producing Permits using Claims provided
func NewAuthMethod(method string, claim Claim) (AuthMethod, error) {
	switch method {
	case "secret":
		return NewAuthMethodSecret(claim)
	}

	return nil, fmt.Errorf("Invalid method requested: %s", method)
}
