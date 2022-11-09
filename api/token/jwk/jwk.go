// package jwk implements utilities related to the maintenance of JWKs
//
// JSON Web Key (JWK) - (RFC 7517) are used to sign/verify JWT signature,
// keys are provided by the issuer's API, and rarely change.
// MS (Azure) recommends updating them every 24h, and on demand if new keys
// are used.
// The implemented JWK cache will be updated if older than acceptable age.
// Updates will be forced on demand if keys are not found in current Cache.
// Updater prevents updates more frequent than allowed.,
// ToDo:
//   - PERF: cache cleanup.  Current implementation only updates and adds to cache.
//   - SEC: source API cert verification.  Currently depends on OS Cert Store
//   - package is made to handle MS key store, other JWKs may have different ways of doing it.
//     Consider compatibility with other IDPs.
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
	"sync"
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

// structure for the in-memory cache of current JWKs
// jwkMap - key: kid; value: JWK
type JwkSetCache struct {
	jwkMap             map[string]JWK
	lastUpdatedAt      int64
	apiUri             string
	apiTimeoutMs       int
	maxRefreshInterval int
	minRefreshInterval int
	mu                 sync.Mutex
}

// collecion of keys, as returned by MS Azure API
type JwkSet struct {
	Keys []JWK `json:"keys"`
}

// create new JWK Cache - used by unit tests to init JWK Cache for various scenarios.
// Normally JWKSet Cache is is initialized by config.go using jwk.InitJWKCache()
// CONSIDER:  creation and initialization should perhaps be done in config, not here.
func NewJWKCache(apiUri string, apiTimeoutMs, maxRefreshInterval, minRefreshInterval int) *JwkSetCache {

	jsc := &JwkSetCache{
		jwkMap:             make(map[string]JWK),
		apiUri:             apiUri,
		apiTimeoutMs:       apiTimeoutMs,
		maxRefreshInterval: maxRefreshInterval,
		minRefreshInterval: minRefreshInterval,
	}
	log.Infof("JWK Cache crated for %v", apiUri)
	return jsc
}

var JWKSetCache JwkSetCache

// Initialize JWK Cache
// recommended values for azure: timeout: 1000ms
// maxRefreshInterval: 300s, minRefreshInterval: 86400
// AZURE API: "https://login.microsoftonline.com/common/discovery/v2.0/keys"
func InitJWKCache() {
	//TODO: move JWK Cache settings to global config
	apiUri := "https://login.microsoftonline.com/common/discovery/v2.0/keys"
	maxRefreshIntervalSec := 300
	minRefreshIntervalSec := 86400
	apiTimeoutMs := 1000

	JWKSetCache = JwkSetCache{
		jwkMap:             make(map[string]JWK),
		apiUri:             apiUri,
		apiTimeoutMs:       apiTimeoutMs,
		maxRefreshInterval: maxRefreshIntervalSec,
		minRefreshInterval: minRefreshIntervalSec,
	}
	log.Infof("JWK Cache crated for %v", apiUri)

	err := JWKSetCache.Update()
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
func (j *JwkSet) updateFromSource(uri string, timeoutMils int) error {

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

// update JwkSetCache with contents from the source API (i.e. microsoft)
//
// Custom errors:
// "JWK source error: nothing received, cache not updated"
func (j *JwkSetCache) Update() error {

	//block parallel attempts to update cache, one update is enough
	j.mu.Lock()
	defer j.mu.Unlock()

	//verify last update timestamp vs. max update interval
	//this is to prevent high frequency updates - respect the source api!
	if j.lastUpdatedAt >= time.Now().Unix()-int64(j.maxRefreshInterval) {
		log.Infof("JWK Cache max refresh interval exceeded, current cache is %vs old.", time.Now().Unix()-j.lastUpdatedAt)
		return nil
	}

	var newJwkSet JwkSet

	// get new JWKs from source API
	log.Debugf("JWK: getting new keys from MS API.")
	err := newJwkSet.updateFromSource(j.apiUri, j.apiTimeoutMs)
	if err != nil {
		log.Errorf(fmt.Sprint(err))
		return errors.New("JWK: source error: nothing received, cache not updated")
	}

	// add new keys from source into JWKSetCache
	newKeyCounter := 0
	for _, sourceKey := range newJwkSet.Keys {
		if _, ok := j.jwkMap[sourceKey.Kid]; !ok {
			//new keys, update map
			j.jwkMap[sourceKey.Kid] = sourceKey
			newKeyCounter++
		}
	}
	j.lastUpdatedAt = time.Now().Unix()
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
func (j *JwkSetCache) x509cert(kid string) (string, error) {

	jwk, ok := j.jwkMap[kid]
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
func (j *JwkSetCache) RsaPubKey(kid string) (*rsa.PublicKey, error) {

	//verify last update timestamp vs. min update interval
	//this is to force an update if cache is too stale.
	if j.lastUpdatedAt < time.Now().Unix()-int64(j.minRefreshInterval) {
		log.Infof("JWK Cache is stale, forcing update.")
		err := j.Update()
		if err != nil {
			log.Warnf("Problems updating JWK Cache, cache is stale. [%vs]",
				time.Now().Unix()-j.lastUpdatedAt)
		}

	}

	//check if kid exists in current cache and force an update if kid not found
	//and cache is not fresh
	_, ok := j.jwkMap[kid]
	if !ok && j.lastUpdatedAt < time.Now().Unix()-int64(j.maxRefreshInterval) {
		log.Infof("kid [%v] not found in JWK Cache, updating cache", kid)
		err := j.Update()
		if err != nil {
			return nil, errors.New("kid not found in JWKSetCache")
		}
	}

	var err error

	//retrieve PEM x.509 Cert from JWKSetCache
	var pemCert string
	pemCert, err = j.x509cert(kid)
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
