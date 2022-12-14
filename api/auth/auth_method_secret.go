package auth

import (
	"fmt"
	"context"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	log "github.com/sirupsen/logrus"

	"admincheckapi/api/auth/tokencache"	
)

// MethodSecret is a container for Claims and Permits obtained
type AuthMethodSecret struct {
	Claim
	Permit
}

// Vanilla unsafe cache, full implementation TBD
var (
	cacheAccessor = &tokencache.TokenCache{"cache.json"}
)

// acquireTokenClientSecret does auth request for a method
func acquireTokenClientSecret(claim Claim) (string, error) {
	crd, err := confidential.NewCredFromSecret(claim.ClientSecret)
	if err != nil {
		return "", fmt.Errorf("Error creating credentials with secret")
	}

	app, err := confidential.New(claim.ClientID,
		crd,
		confidential.WithAuthority(claim.Authority),
		confidential.WithAccessor(cacheAccessor))
	if err != nil {
		return "", fmt.Errorf("Error creating confidential")
	}

	result, err := app.AcquireTokenSilent(context.Background(), claim.Scopes)
	if err != nil {
		result, err = app.AcquireTokenByCredential(context.Background(),
			claim.Scopes)
		if err != nil {
			return "", fmt.Errorf("Error acquire tocken with credential")
		}
	}

	return result.AccessToken, nil
}

// NewAuthMethodSecret creates new object with original claim and a permit
func NewAuthMethodSecret(claim Claim) (AuthMethodSecret, error) {
	log.Debugf("Requested auth with claim: %+v", claim)
	token, err := acquireTokenClientSecret(claim)
	if err != nil {
		return AuthMethodSecret{
			Claim{},
			Permit{},
		},
			fmt.Errorf("Error requesting token with secret claim: %+v", claim)
	}

	return AuthMethodSecret{claim, Permit{token}}, nil
}

// Token gives out the permit artefact
func (m AuthMethodSecret) Token() string {
	return m.Permit.token
}
