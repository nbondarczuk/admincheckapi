package awssm

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	log "github.com/sirupsen/logrus"
)

//
// GetSecret fetches the key -> val mapping of tenant id to MSAD login credentials
//
func (ss *AWSSecretStorage) GetSecret(secretName string) (string, error) {
	// New WAS session
	log.Debugf("Opening new AWS session")
	s, err := session.NewSession()
	if err != nil {
		return "", err
	}
	log.Debugf("AWS session opened")

	// Get hold of AWS secret manager
	log.Debugf("Accessing AWS secret manager")
	svc := secretsmanager.New(s, aws.NewConfig().WithRegion(ss.Region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}
	log.Debugf("AWS secret manager connected")

	// Get the secrets from aWS secret manager
	log.Debugf("Getting AWS secret: %s", secretName)
	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", fmt.Errorf("Error while gettting secret: %s %s", secretName, err.Error())
	}
	if result.SecretString != nil {
		log.Debugf("Got AWS secret")
		return *result.SecretString, nil
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			return "", err
		}
		return string(decodedBinarySecretBytes[:len]), nil
	}
	
	return "", nil
}
