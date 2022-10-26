// package jwk implements utilities related to the maintenance of JWKs
//
// JSON Web Key (JWK) - (RFC 7517) are used to sign/verify JWT signature,
// keys are provided by the issuer's API, and rarely change.
// MS (Azure) recommends updating them every 24h, and on demand if new keys
// are used.
//
// ToDo:
// * SEC: source API cert verification.  Currently depends on OS Cert Store
// * FUNC: auto update if last update date > 24h ago - perhaps upon first kid not found req.
// * FUNC/SEC: update on demand (if last update > 5min ago, to protect against flood attacks)
// * PERF: hold ongoing (get) RsaPubKey requests if JWK Set is being updated, at least for 2x
//   API Timeout
// * package is made to handle MS key store, other JWKs may have different ways of doing it.
//   Consider compatibility with other IDPs.

package jwk

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// JSON Web Key (JWK) - (RFC 7517)
// only fields used by Azure are used!
// filed X5c is a slice (list of elements) but in the actual JWK Set issued by
// Azure it's just one value - the rsa certificate.
// in case there where more, it's a cert chain, the first one is used for JWT
// signing.
// Kid - Key Id used in JWT's to indicate public key to be used to verification
type JWK struct {
	Kty    string   `json:"kty"`
	Use    string   `json:"use"`
	Kid    string   `json:"kid"`
	X5t    string   `json:"x5t"`
	N      string   `json:"n"`
	E      string   `json:"e"`
	X5c    []string `json:"x5c"`
	Issuer string   `json:"issuer"`
}

// in-memory cache of current JWKs
// Map - key: kid; value: JWK
type JwkSetMap map[string]JWK

// collecion of keys, as returned by MS Azure API
type JwkSet struct {
	Keys []JWK `json:"keys"`
}

var JWKSetCache = make(JwkSetMap)

func InitJWKCache() {
	msJwkSourceApi := "https://login.microsoftonline.com/common/discovery/v2.0/keys"
	err := JWKSetCache.Update(msJwkSourceApi, 1000)
	if err != nil {
		log.Fatalf("Error loading JWK Set Cache [%s]", fmt.Sprint(err))
	}
}

// Acquire current JWK Set from <uri> API and populate JwkSet struct with JWKs.
// Params:
// uri: address of the source API
// timeouteMils: timeout, reasonable: 500ms-1000ms
// Azure key store:
// https://login.microsoftonline.com/common/discovery/v2.0/keys
// !TODO: this URI should be a part of standard configuration
// !CONSIDER: consider moving timeout to global config
func (j *JwkSet) UpdateFromSource(uri string, timeoutMils int) error {

	c := http.Client{Timeout: time.Duration(timeoutMils) * time.Millisecond}
	resp, err := c.Get(uri)
	if err != nil {
		log.Errorf("JWK Source API connection error: [%v]", err)
		return fmt.Errorf("JWK Source API connection error: [%v]", err)
	}
	defer resp.Body.Close()

	// continue only if API returned HTTP Status 200 OK
	if resp.StatusCode != 200 {
		log.Errorf("JWK Source API error: [%v]", resp.Status)
		return errors.New("JWK Source API error: " + resp.Status)
	}

	JwkSetJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(JwkSetJSON), &j)
	if err != nil {
		return err
	}
	return nil
}

// update JwkSetMap with contents from the source API (i.e. microsoft)
// Params:
// uri - source JWK API address
// timeoutMils - timeout in milliseconds (suggested > 500ms)
// Custom errors:
// "JWK source error: nothing received, cache not updated"
func (j JwkSetMap) Update(uri string, timeoutMils int) error {

	var newJwkSet JwkSet

	// get new JWKs from source API
	log.Debugf("JWK: getting new keys from API.")
	err := newJwkSet.UpdateFromSource(uri, timeoutMils)
	if err != nil {
		log.Errorf(fmt.Sprint(err))
		return errors.New("JWK: source error: nothing received, cache not updated")
	}

	// add new keys from source into JWKSetCache
	newKeyCounter := 0
	for _, sourceKey := range newJwkSet.Keys {
		if _, ok := j[sourceKey.Kid]; !ok {
			//new keys, update map
			j[sourceKey.Kid] = sourceKey
			newKeyCounter++
		}
	}
	if newKeyCounter > 0 {
		log.Debugf("JWK: JWKSetCache updated with %v new JWKs.", newKeyCounter)
	} else {
		log.Debugf("JWK: JWKSetCache not updated, now new JWKs found at source API")
	}
	return nil
}

// x509 cert for a given key id (kid) in JwkSetMap
//
// custom errors:
// "kid not found in JWKSetCache"
// "certificate filed empty for given kid"
func (j JwkSetMap) X509cert(kid string) (string, error) {

	jwk, ok := j[kid]
	if !ok {
		return "", errors.New("kid not found in JWKSetCache")
	}

	cert := jwk.X5c[0]
	// CONSIDER: X5c field might be empty for some reason - is this an error?
	if len(cert) == 0 {
		log.Errorf("JWK: certificate filed empty for kid: " + kid)
		return "", errors.New("JWK certificate filed empty for given kid")
	}

	//add text tags to cert string as per x509 spec
	if !strings.Contains(cert, "-----BEGIN CERTIFICATE-----") {
		log.Debugf("adding text tags to bare certificate string")
		cert = "-----BEGIN CERTIFICATE-----\n" + cert + "\n-----END CERTIFICATE-----"
	}
	return cert, nil
}

// RSA Public Key for a given key id (kid) in JwkSetMap
func (j JwkSetMap) RsaPubKey(kid string) (*rsa.PublicKey, error) {
	var err error

	//retrieve PEM x.509 Cert from JWKSetCache
	var pemCert string
	pemCert, err = j.X509cert(kid)
	if err != nil {
		return nil, err
	}

	//extract pub key from certificate
	var rsaPubKey *rsa.PublicKey
	rsaPubKey, err = pubKeyfromX509Cert(pemCert)
	if err != nil {
		return nil, err
	}
	return rsaPubKey, nil
}

// Extracts RSA public key from a giver x.509 certificate
// Custom errors:
// "failed to decode PEM block containing certificate from certString"
func pubKeyfromX509Cert(certPEMString string) (*rsa.PublicKey, error) {
	var err error

	//decode cert data from certString
	block, _ := pem.Decode([]byte(certPEMString))
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("failed to decode PEM block containing certificate from certString")
	}

	var cert *x509.Certificate
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	//extract public key from cert and return
	return cert.PublicKey.(*rsa.PublicKey), nil
}
