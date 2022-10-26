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
	ClientID            string
	Authority           string
	Scopes              []string
	Username            string
	Password            string
	RedirectURI         string
	CodeChallenge       string
	CodeChallengeMethod string
	State               string
	ClientSecret        string
	Thumbprint          string
	PemData             string
}

// NewAuthMethod is a factory producing Permits using Claims provided
func NewAuthMethod(method string, claim Claim) (AuthMethod, error) {
	switch method {
	case "secret":
		return NewAuthMethodSecret(claim)
	}

	return nil, fmt.Errorf("Invalid method requested: %s", method)
}
