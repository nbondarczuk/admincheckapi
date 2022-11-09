package jwk_test

import (
	"admincheckapi/api/token/jwk"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitJWKCache(t *testing.T) {

	jwk.InitJWKCache()

	err := jwk.JWKSetCache.Update()
	assert.NoError(t, err, "expecting no error")
}

func TestJWKCacheFrequentUpdates(t *testing.T) {
	var JWKSetCache = jwk.NewJWKCache("https://login.microsoftonline.com/common/discovery/v2.0/keys",
		1000, 2, 86400)

	err := JWKSetCache.Update()
	time.Sleep(time.Second * 1)
	err = JWKSetCache.Update()

	assert.NoError(t, err, "expecting no error")
}

// func TestInitJWKCache(t *testing.T) {
// 	jwk.InitJWKCache()
// 	assert.Equal(t, true, true)
// }

// Case: negative, expected: connection error
func TestUpdateJWKCache_ConErr(t *testing.T) {
	var JWKSetCache = jwk.NewJWKCache("https://example-asdf.com/nonexistant",
		1000, 300, 86400)
	err := JWKSetCache.Update()
	assert.ErrorContains(t, err, "JWK: source error:", "expecting source error")
}

// Case: negative, expected: 404 NotFound
func TestUpdateJWKCache_404Err(t *testing.T) {
	var JWKSetCache = jwk.NewJWKCache("https://login.microsoftonline.com/common/discovery/v2.0/keys123",
		1000, 300, 86400)
	err := JWKSetCache.Update()
	assert.ErrorContains(t, err, "JWK: source error:", "expecting source error")
}

// Case: negative, expected timeout  (giving it 10ms to timeout)
func TestUpdateJWKCache_TimeoutErr(t *testing.T) {
	var JWKSetCache = jwk.NewJWKCache("https://login.microsoftonline.com/common/discovery/v2.0/keys123",
		50, 300, 86400)
	err := JWKSetCache.Update()
	assert.ErrorContains(t, err, "JWK: source error:", "expecting source error")
}

func TestJwkCacheRsaPubKey(t *testing.T) {

	//positive test, kid + key exist.
	testkid := "2ZQpJ3UpbjAYXYGaXEJl8lV0TOI"
	pubKeyExpStr := "24270816892089731719543783751432048815491545991257707613897766129285086092017663113913065374578252866292844529129303576012291981496377741717415683430954480719297040237552883183880720525825059504957853986049056158452751101633284726834183143745634822061761330491404604872173147765504606958311929588316683430193995492355659914947719650469297035699306817588668110003001642212415645738806065451691699269723426112951520848277600352153453283413356059207174622264137605364529710378358949775695884543447084486374661976470988089538432223651073995367128692155111525981268505073899618934526964112047508332751054877035214512639273"

	var JWKSetCache = jwk.NewJWKCache("https://login.microsoftonline.com/common/discovery/v2.0/keys",
		1000, 3, 5)
	err := JWKSetCache.Update()

	pubKey, err := JWKSetCache.RsaPubKey(testkid)
	assert.Equal(t, fmt.Sprint(pubKey.N), pubKeyExpStr, "should find a key in JWK Set Cache")
	assert.NoError(t, err, "expecting no errors")

	//negative test, kid not in JWKSetCache
	_, err = JWKSetCache.RsaPubKey("noSuchKid")
	assert.ErrorContains(t, err, "kid not found", "nonexistent kid found in JWKSetCache")

	time.Sleep(time.Second * 6)
	pubKey, err = JWKSetCache.RsaPubKey(testkid)
	assert.NoError(t, err, "expecting no errors")

}

func TestJwkCacheRsaPubKeyWithUpdate(t *testing.T) {

	//positive test, kid + key exist.
	testkid := "noSuchKid"

	var JWKSetCache = jwk.NewJWKCache("https://login.microsoftonline.com/common/discovery/v2.0/keys",
		1000, 3, 10)
	err := JWKSetCache.Update()

	time.Sleep(time.Second * 4)
	_, err = JWKSetCache.RsaPubKey(testkid)
	//negative test, kid not in JWKSetCache
	_, err = JWKSetCache.RsaPubKey("noSuchKid")
	assert.ErrorContains(t, err, "kid not found", "nonexistent kid found in JWKSetCache")

}
