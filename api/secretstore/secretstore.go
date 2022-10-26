package secretstore

import (
	"fmt"

	"admincheckapi/api/auth"
	"admincheckapi/api/config"
)

const MSOnlineURL = "https://login.microsoftonline.com"

var JWTSecretToken string

//
// 1 item cache - shall be replaced with a hit to AWS safe key storage
//
func TenantJWTToken() (string, error) {
	if JWTSecretToken != "" {
		return JWTSecretToken, nil
	}

	var claims auth.Claim = auth.Claim{
		Authority:    fmt.Sprintf("%s/%s", MSOnlineURL, config.Setup.TenantId),
		ClientID:     config.Setup.ClientId,
		Scopes:       []string{".default"},
		ClientSecret: config.Setup.ClientSecret,
	}

	am, err := auth.NewAuthMethod("secret", claims)
	if err != nil {
		return "", fmt.Errorf("Error while creating %s (%+v) auth method: %s", "secret", claims, err)
	}

	JWTSecretToken = am.Token()

	return JWTSecretToken, nil
}
