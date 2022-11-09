package secretstore

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"admincheckapi/api/auth"
	"admincheckapi/api/aws/awssm"
	"admincheckapi/api/config"
)

const MSOnlineURLPattern = "https://login.microsoftonline.com/%s"

var JWTSecretToken string

type CredentialsSecret struct {
	Authority    string
	ClientID     string
	Scopes       []string
	ClientSecret string
}

//
// AWS secret store mapping the token id to MS graph crdentials needed
// to obtain an access token
//
func TenantJWTToken(tenantId string) (string, error) {
	log.Tracef("Begin: TenantJWTToken")
	defer log.Tracef("End: TenantJWTToken")

	// Try using local lru cache. Cache refresh is TBD
	if JWTSecretToken != "" {
		log.Debugf("Got jwt token from cache")
		return JWTSecretToken, nil
	}

	var credsSecret CredentialsSecret
		
	if config.Setup.AWSUseSecretStore {
		// Get AWS secret store handle
		region := os.Getenv("AWS_REGION")
		var ss awssm.AWSSecretStorage
		if region != "" {
			ss = awssm.AWSSecretStorage{Region: region}
		}
		log.Debugf("Using AWS secret store region: %s", region)
		
		// The secret may be labelled in a flexible way as AWS ecrets are inmutable
		creds, err := ss.GetSecret(config.Setup.SecretNamePrefix + tenantId)
		if err != nil {
			return "", fmt.Errorf("Error while getting secret from secret storage: %v", err)
		}
		log.Debugf("Got from AWS secret store credentials for MS graph token: %+v", creds)
		
		// Decode the resulting json with secret store data
		var dataFromSS map[string]string
		err = json.Unmarshal([]byte(creds), &dataFromSS)
		if err != nil {
			return "", fmt.Errorf("Error while decoding secret from storage: %v", err)
		}
		log.Debugf("Decoded data from secret store: %+v", dataFromSS)
		
		// Get the secrets from secret store data
		err = json.Unmarshal([]byte(dataFromSS[tenantId]), &credsSecret)
		if err != nil {
			return "", fmt.Errorf("Error while decoding secret from credentials: %v", err)
		}
		log.Debugf("Decoded credentials from secret store: %+v", credsSecret)
	} else {
		log.Debugf("Using secrets from env: %s")
		credsSecret.Authority = config.Setup.Authority
		credsSecret.ClientID  = config.Setup.ClientId
		credsSecret.Scopes = config.Setup.Scopes
		credsSecret.ClientSecret = config.Setup.ClientSecret
	}
	
	// Connecting to MS graph to get a service token
	var claims auth.Claim = auth.Claim{
		Authority:    credsSecret.Authority,
		ClientID:     credsSecret.ClientID,
		Scopes:       credsSecret.Scopes,
		ClientSecret: credsSecret.ClientSecret,
	}
	am, err := auth.NewAuthMethod("secret", claims)
	if err != nil {
		return "", fmt.Errorf("Error while creating %s (%+v) auth method: %s",
			"secret", claims, err)
	}
	log.Debugf("Connected to MS graph with auth method: %s claims: %v", "secret", claims)

	// Parse the new service token
	JWTSecretToken = am.Token()
	log.Debugf("Got new jwt token from MS Azure")

	return JWTSecretToken, nil
}
