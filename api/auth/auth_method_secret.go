package auth

import (
	"context"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"

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
		return "", err
	}

	app, err := confidential.New(claim.ClientID,
		crd,
		confidential.WithAuthority(claim.Authority),
		confidential.WithAccessor(cacheAccessor))
	if err != nil {
		return "", err
	}

	result, err := app.AcquireTokenSilent(context.Background(), claim.Scopes)
	if err != nil {
		result, err = app.AcquireTokenByCredential(context.Background(),
			claim.Scopes)
		if err != nil {
			return "", err
		}
	}

	return result.AccessToken, nil
}

// NewAuthMethodSecret creates new object with original claim and a permit
func NewAuthMethodSecret(claim Claim) (AuthMethodSecret, error) {
	token, err := acquireTokenClientSecret(claim)
	if err != nil {
		return AuthMethodSecret{Claim{}, Permit{}}, err
	}

	return AuthMethodSecret{claim, Permit{token}}, nil
}

// Token gives out the permit artefact
func (m AuthMethodSecret) Token() string {
	return m.Permit.token
}
