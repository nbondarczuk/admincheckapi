package jwk_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"admincheckapi/api/token/jwk"
)

// initialize JWKSetCache
var JWKSetCache = make(jwk.JwkSetMap)

func TestGetJWKsFromSource(t *testing.T) {

	// positive tests

	// expected non-zero length JwkSet (ie. some keys were returned, typically 8)
	// testing against actuall API.  behave with respect!!!
	testURI := "https://login.microsoftonline.com/common/discovery/v2.0/keys"
	var newJwkSet jwk.JwkSet
	err := newJwkSet.UpdateFromSource(testURI, 1000)
	assert.GreaterOrEqual(t, len(newJwkSet.Keys), 1, "loading of JWKs from source failed")

	// negative tests

	// expected: connection error
	testURI = "https://example-asdf.com/nonexistant"
	err = newJwkSet.UpdateFromSource(testURI, 1000)
	assert.ErrorContains(t, err, "API connection error", "TEST: expecting API connection error")

	// expected: 404 NotFound
	testURI = "https://login.microsoftonline.com/common/discovery/v2.0/keys123"
	err = newJwkSet.UpdateFromSource(testURI, 1000)
	assert.EqualError(t, err, "JWK Source API error: 404 Not Found", "TEST: expecting 404 error")

	// expected timeout  (giving it 10ms to timeout)
	testURI = "https://example.com"
	err = newJwkSet.UpdateFromSource(testURI, 10)
	assert.ErrorContains(t, err, "Timeout", "Timeout error")
}

func TestUpdateJWKSetCache(t *testing.T) {

	// positive tests

	// expected: no errors, length of JWKSetCache >= then before test.
	// testing against actual API, expecting ~9JWKs behave with respect!!!
	testURI := "https://login.microsoftonline.com/common/discovery/v2.0/keys"
	JwkSetCacheLenBeforeTest := len(JWKSetCache)
	err := JWKSetCache.Update(testURI, 1000)

	assert.GreaterOrEqual(t, len(JWKSetCache), JwkSetCacheLenBeforeTest, err)
	assert.NoError(t, err, err)
}

func TestJwkSetMapX509cert(t *testing.T) {

	//positive test, kid + cert exist
	testkid := "2ZQpJ3UpbjAYXYGaXEJl8lV0TOI"
	cert, err := JWKSetCache.X509cert(testkid)
	assert.GreaterOrEqual(t, len(cert), 100, err)
	assert.NoError(t, err, err)

	//negative test, kid not in JWKSetCache
	cert, err = JWKSetCache.X509cert("noSuchKid")

	assert.ErrorContains(t, err, "kid not found", "test error: nonexistent kid found in JWKSetCache")

}

func TestJwkSetMapRsaPubKey(t *testing.T) {

	//positive test, kid + key exist
	testkid := "2ZQpJ3UpbjAYXYGaXEJl8lV0TOI"
	_, err := JWKSetCache.RsaPubKey(testkid)
	//assert.GreaterOrEqual(t, len(cert), 100, err)
	assert.NoError(t, err, err)

	//negative test, kid not in JWKSetCache
	_, err = JWKSetCache.RsaPubKey("noSuchKey")

	assert.ErrorContains(t, err, "kid not found", "test error: nonexistent kid found in JWKSetCache")

}
