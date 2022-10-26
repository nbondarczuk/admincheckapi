
// package jwk implements utilities related validating and parsing JWTs
//
// JWTs signed with RS256 signatures are validated against
// JSON Web Keys (JWK) maintained by the jwk package initialized during
// main config.
// other checks are performed as well mainly to leave detailed trace
// in logs for security warnings and audit.

package token

import (
	"fmt"
	"regexp"

	jwt "github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"admincheckapi/api/token/jwk"
)

// Token is a container for original token and interesting claims
type (
	Token struct {
		TokenPayload []byte
		Groups       []string
	}

	msTokenClaims struct {
		jwt.StandardClaims
		Groups []string `json:"groups,omitempty"`
	}
)

// NewToken parses the JWT token structure
func NewToken(payload []byte) (Token, error) {
	claims, err := parseJWTClaims(string(payload))
	if err != nil {
		return Token{}, err
	}

	return Token{TokenPayload: payload, Groups: claims.Groups}, nil
}

// AdminGroup extracts roles as group fields of the JSON encoded JWT token - first one only
func (t Token) AdminGroup() (string, error) {
	if len(t.Groups) == 0 {
		return "", fmt.Errorf("No admin group found in the token")
	}

	return t.Groups[0], nil
}

// AdminGroups gets all admin groups found in the token
func (t Token) AdminGroups() ([]string, error) {
	return t.Groups, nil
}

// parse JWT token, attempt sig verification vs. public key from JwkSetCache, returns token claims and err if any
// consider error type before deciding to use unverified token.
func parseJWTClaims(tokenStr string) (*msTokenClaims, error) {

	//verify token format, just in case something weird makes it here.
	// expecting base64_chars.base64chars.signature_chars
	if tokenFormatMatch, _ := regexp.MatchString(`^([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_=]+)\.([a-zA-Z0-9_\-\+\/=]*)$`, tokenStr); !tokenFormatMatch {
		return nil, fmt.Errorf("token format error")
	}

	var tokenClaims msTokenClaims

	token, err := jwt.ParseWithClaims(tokenStr, &tokenClaims, func(token *jwt.Token) (interface{}, error) {

		//verify signature alg. used in token, if different from RS256, parse without signature verification
		if alg := token.Header["alg"]; fmt.Sprint(alg) != "RS256" {
			log.Infof("Unexpected signing method: [%s], attempting to parse without sig verification", alg)
			return nil, fmt.Errorf("unexpected signing method: %v", alg)
		}

		//retrieve key id from header, then public key from JWK Set Cache.
		if kid, ok := token.Header["kid"]; !ok {
			log.Infof("kid not found in token header, attempting to parse without sig verification")
			return nil, fmt.Errorf("token does not contain key id (Header.kid), unable to verify signature")
		} else {
			rsaPublicKey, err := jwk.JWKSetCache.RsaPubKey(fmt.Sprint(kid))
			if err != nil {
				log.Infof("public key not found in cache, attempting to parse without sig verification, kid: [%s]", fmt.Sprint(kid))
				return nil, fmt.Errorf("no public key found for kid: [%s]", kid)
			}
			log.Infof("public key found, attempting to parse with sig verification")
			return rsaPublicKey, nil
		}
	})

	if err != nil {
		// parsing errors encountered, check if token is a "JWT" type
		if fmt.Sprint(token.Header["typ"]) == "JWT" {
			// log the token for potential security analysis.
			// Return &tokenClaims and let the caller decide what to do with it.
			log.Warnf("JWT parsed, found issues: [%s], token: [%s]", err.Error(), token.Raw)
			return &tokenClaims, err
		} else {
			log.Infof("JWT parsing error: [%s]", err.Error())
			return nil, err
		}
	}

	// all good, token is healthy
	return &tokenClaims, nil
}
