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

	"admincheckapi/api/token/jwk"

	jwt "github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

// Token is a container for original token and interesting claims
// so far tid (tokenId) and groups are of interest
type (
	Token struct {
		TokenPayload []byte
		Groups       []string
		Tid          string
		Oid          string
		Idtyp        string
	}

	msTokenClaims struct {
		jwt.StandardClaims
		Groups []string `json:"groups,omitempty"`
		Tid    string   `json:"tid,omitempty"`
		Oid    string   `json:"oid,omitempty"`
		Idtyp  string   `json:"idtyp,omitempty"`
	}
)

// NewToken parses the JWT token structure
func NewToken(payload []byte) (Token, error) {
	claims, err := parseJWTClaims(string(payload))
	if err != nil {
		return Token{}, err
	}

	return Token{TokenPayload: payload,
		Groups: claims.Groups,
		Tid:    claims.Tid,
		Oid:    claims.Oid,
		Idtyp:  claims.Idtyp}, nil
}

// AdminGroups gets all admin groups found in the token
func (t Token) AdminGroups() ([]string, error) {
	return t.Groups, nil
}

// TenantId is needed to log into client's organization
func (t Token) TenantId() (string, error) {
	return t.Tid, nil
}

// Object Id in the Microsoft identity system, in this case, a user account or an application.
// In MS Tokens it's the "oid" field.
// Returns String, a GUID
func (t Token) ObjectId() (string, error) {
	return t.Oid, nil
}

// Check if the token was issued to an application of a user.
// idtyp is an optional claim in the Microsoft identity system - used to distinguish between app-only
// access tokens and access tokens for users.
// For applicaions idtyp = "app", for users is not included.
func (t Token) IsApp() (bool, error) {
	if t.Idtyp == "app" {
		return true, nil
	}
	return false, nil
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
			log.Debugf("public key found, attempting to parse with sig verification")
			return rsaPublicKey, nil
		}
	})

	if err != nil {
		// parsing errors encountered, check if token is a "JWT" type
		if fmt.Sprint(token.Header["typ"]) == "JWT" {
			// log the token for potential security analysis.
			// Return &tokenClaims and let the caller decide what to do with it.
			log.Warnf("JWT parsed, found issues: [%s], token: [%s]", err.Error(), token.Raw)
			return &tokenClaims, nil
		} else {
			log.Infof("JWT parsing error: [%s]", err.Error())
			return nil, err
		}
	}

	// all good, token is healthy
	return &tokenClaims, nil
}
