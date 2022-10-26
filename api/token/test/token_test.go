package token_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"admincheckapi/api/token"
	"admincheckapi/api/token/jwk"
	"admincheckapi/test/testconfig"

	"github.com/stretchr/testify/assert"
)

type jwtParcingTest struct {
	testDesc    string
	tokenStr    string
	expectedErr string
	gropuCount  int
}

var jwtParcingTests = []jwtParcingTest{

	//various test cases with diffent tokens and expected errors:
	jwtParcingTest{
		"POS: sig valid, not expired, 3 groups - FRESH TOEKN EXPECTED!",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSIsImtpZCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSJ9.eyJhdWQiOiJhcGk6Ly8wNDJkODA3ZC1hMThiLTQ5NjUtYTgyNC1lZmY3Mzg1NjA3ZTYiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC8yM2Q0MjUwZS0zNzExLTQ4ZTEtOGRmZS1kMjdmMjg1MTU0YTEvIiwiaWF0IjoxNjY2MTAwNDk5LCJuYmYiOjE2NjYxMDA0OTksImV4cCI6MTY2NjEwNDM5OSwiYWlvIjoiRTJaZ1lHQ1ZMcGdTWlJGdXE1WHI1SzZocytjZEFBPT0iLCJhcHBpZCI6IjUzNTllNmEzLWRjMWUtNGFjMS04MDczLWQwZTY1OGEwMDJjNCIsImFwcGlkYWNyIjoiMSIsImdyb3VwcyI6WyI0YzllMzhjZS1hZjY3LTRjMDgtYmRlZS1iZThmMDZiZmFjNmIiLCIyZTEyZDFiYi1jNDhlLTQyY2QtYjRjZC01NzMzODVlNGE5ZGMiLCJjNzExM2EzNi1kZTRkLTRmY2MtOWVhMC01MjM3MjcxMGQzY2QiXSwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvMjNkNDI1MGUtMzcxMS00OGUxLThkZmUtZDI3ZjI4NTE1NGExLyIsIm9pZCI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInJoIjoiMC5BWUlBRGlYVUl4RTM0VWlOX3RKX0tGRlVvWDJBTFFTTG9XVkpxQ1R2OXpoV0ItYVZBQUEuIiwicm9sZXMiOlsiYXJnb25BZG1pblJvbGUxIiwidGVzdFJvbGUxIl0sInN1YiI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInRpZCI6IjIzZDQyNTBlLTM3MTEtNDhlMS04ZGZlLWQyN2YyODUxNTRhMSIsInV0aSI6ImFMZlZ1ajVYMzBLWUdtNzFWMXRRQUEiLCJ2ZXIiOiIxLjAifQ.l-6eTGXOmPlCXxVWH2Z4gYAU9t2ClWv8YtvfIVgEmnget_PTXKuBQmXBLfW3j7NiSR7OoGmfCjyy06TXccGNjmIHNyevXq4_biBfqbaiFEbpPDYO02w57dJr0vUrNaYJwnKUKRLyM3jC2mTniYHJtRqBwaXdHkX-bknKt26XNnyddOqumcsxFE-brmYgQMcpJUxMFnBzoM8VOmAsThJkw1JYZn2hnQ-eOSNu5yGLnZKw1a-rWKfZc9_rJe_xxaUaQGb4p_NmpZXtEAjd7RJWDmps2H0kLY76VQ3RPsRH-JR1K4sxE5aeCCKMGXxAejmLfo5z5aIyArQDY-dpEqzTXw",
		"", // this is the positive test, new fresh token is expected.
		3},

	jwtParcingTest{
		"NEG: sig valid, expired, 3 groups",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSIsImtpZCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSJ9.eyJhdWQiOiJhcGk6Ly8wNDJkODA3ZC1hMThiLTQ5NjUtYTgyNC1lZmY3Mzg1NjA3ZTYiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC8yM2Q0MjUwZS0zNzExLTQ4ZTEtOGRmZS1kMjdmMjg1MTU0YTEvIiwiaWF0IjoxNjY1OTk5MzI2LCJuYmYiOjE2NjU5OTkzMjYsImV4cCI6MTY2NjAwMzIyNiwiYWlvIjoiRTJaZ1lLaStuWnA3NmQ0TDdRc0Y2Zzk1TGF3NUFBPT0iLCJhcHBpZCI6IjUzNTllNmEzLWRjMWUtNGFjMS04MDczLWQwZTY1OGEwMDJjNCIsImFwcGlkYWNyIjoiMSIsImdyb3VwcyI6WyI0YzllMzhjZS1hZjY3LTRjMDgtYmRlZS1iZThmMDZiZmFjNmIiLCIyZTEyZDFiYi1jNDhlLTQyY2QtYjRjZC01NzMzODVlNGE5ZGMiLCJjNzExM2EzNi1kZTRkLTRmY2MtOWVhMC01MjM3MjcxMGQzY2QiXSwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvMjNkNDI1MGUtMzcxMS00OGUxLThkZmUtZDI3ZjI4NTE1NGExLyIsIm9pZCI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInJoIjoiMC5BWUlBRGlYVUl4RTM0VWlOX3RKX0tGRlVvWDJBTFFTTG9XVkpxQ1R2OXpoV0ItYVZBQUEuIiwicm9sZXMiOlsiYXJnb25BZG1pblJvbGUxIiwidGVzdFJvbGUxIl0sInN1YiI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInRpZCI6IjIzZDQyNTBlLTM3MTEtNDhlMS04ZGZlLWQyN2YyODUxNTRhMSIsInV0aSI6Im1SNmhEdkZjWmt1YTdueGdvSkFsQUEiLCJ2ZXIiOiIxLjAifQ.rnPo5jm24js9BSGw6HODQ8V6GwVeWpsU0G9i6CFsyRcJeEPXGLiE5GnTgYpsl40AMIL_uxwRgr_IObt50uIzbdmzZ_9GO_XZLUF_6oH88-UD2M4HNZVkQL3oRluloiv41hFrUPHzTTCxJ0tUKO7YjEuPAeL_To9RQ-Cu_Jdfc5SPwR3KAN1K5numTXZ1szFbI2q0S7eDO9WmhV8XDxO-E9fFI2Bdr2fG8wGQYkJhsPny5of4rt2sREfB1RHDpEHOtDelAp2jfkbiAbP_y4NqDMI8K-7T5jp620Nqj-AUXaTgubTX-V24FWWOtwUkzM606UW-Fm2P_22JXig5Qq35vg",
		"token is expired", // this is the positive test, unless a fresh token is used.
		0},

	jwtParcingTest{
		"NEG: non existing kid, 1 group",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSIsImtpZCI6IlhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWCJ9.eyJhdWQiOiJhcGk6Ly8wNDJkODA3ZC1hMThiLTQ5NjUtYTgyNC1lZmY3Mzg1NjA3ZTYiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC8yM2Q0MjUwZS0zNzExLTQ4ZTEtOGRmZS1kMjdmMjg1MTU0YTEvIiwiaWF0IjoxNjYyOTg1ODA4LCJuYmYiOjE2NjI5ODU4MDgsImV4cCI6MTY2Mjk4OTcwOCwiYWlvIjoiRTJaZ1lPajlmU3lCTlNpeTRldEYyOHhma1RKdkFBPT0iLCJhcHBpZCI6IjUzNTllNmEzLWRjMWUtNGFjMS04MDczLWQwZTY1OGEwMDJjNCIsImFwcGlkYWNyIjoiMSIsImdyb3VwcyI6WyIyZTEyZDFiYi1jNDhlLTQyY2QtYjRjZC01NzMzODVlNGE5ZGMiXSwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvMjNkNDI1MGUtMzcxMS00OGUxLThkZmUtZDI3ZjI4NTE1NGExLyIsIm9pZCI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInJoIjoiMC5BWUlBRGlYVUl4RTM0VWlOX3RKX0tGRlVvWDJBTFFTTG9XVkpxQ1R2OXpoV0ItYVZBQUEuIiwicm9sZXMiOlsiYXJnb25BZG1pblJvbGUxIiwidGVzdFJvbGUxIl0sInN1YiI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInRpZCI6IjIzZDQyNTBlLTM3MTEtNDhlMS04ZGZlLWQyN2YyODUxNTRhMSIsInV0aSI6InU0X1ludC12SjAtVnk2Tk1kWFF4QUEiLCJ2ZXIiOiIxLjAifQ.ZfixXtPMrCEEUrVoW8rUQe-Wqx08nztZ1omqSRmQfWt8dimYndyJ5f4jUkAIfzKvyPtgvpXJo4dvWVpAeNUiPoLR3dcFjrMz9b_EkfXFI3NS0hkSqwFytazHa3v_o_V7TpIi5XBJRgfba5pYJlDeqnEPhEEScOD_jhSTjBDVJWp2j6iEsxmOJ5KKCzyG4FX0laJO16lcGcBsUzkVvIIH2n5FZXyrqMcL2Evgko6d64VAtx4wA5Kxanvs6igS3bIV7MUjpeP67aTcAafVCR20hgWT_IpZ6qaBidr68H6dzyeHTMYfeE-2NvO5o8mrNg65jndIDGA84zyFkxcaqh4ghQ",
		"no public key found for kid",
		0},

	jwtParcingTest{
		"NEG: token with missing kid, 1 group",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSJ9.eyJhdWQiOiJhcGk6Ly8wNDJkODA3ZC1hMThiLTQ5NjUtYTgyNC1lZmY3Mzg1NjA3ZTYiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC8yM2Q0MjUwZS0zNzExLTQ4ZTEtOGRmZS1kMjdmMjg1MTU0YTEvIiwiaWF0IjoxNjYyOTg1ODA4LCJuYmYiOjE2NjI5ODU4MDgsImV4cCI6MTY2Mjk4OTcwOCwiYWlvIjoiRTJaZ1lPajlmU3lCTlNpeTRldEYyOHhma1RKdkFBPT0iLCJhcHBpZCI6IjUzNTllNmEzLWRjMWUtNGFjMS04MDczLWQwZTY1OGEwMDJjNCIsImFwcGlkYWNyIjoiMSIsImdyb3VwcyI6WyIyZTEyZDFiYi1jNDhlLTQyY2QtYjRjZC01NzMzODVlNGE5ZGMiXSwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvMjNkNDI1MGUtMzcxMS00OGUxLThkZmUtZDI3ZjI4NTE1NGExLyIsIm9pZCI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInJoIjoiMC5BWUlBRGlYVUl4RTM0VWlOX3RKX0tGRlVvWDJBTFFTTG9XVkpxQ1R2OXpoV0ItYVZBQUEuIiwicm9sZXMiOlsiYXJnb25BZG1pblJvbGUxIiwidGVzdFJvbGUxIl0sInN1YiI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInRpZCI6IjIzZDQyNTBlLTM3MTEtNDhlMS04ZGZlLWQyN2YyODUxNTRhMSIsInV0aSI6InU0X1ludC12SjAtVnk2Tk1kWFF4QUEiLCJ2ZXIiOiIxLjAifQ.ZfixXtPMrCEEUrVoW8rUQe-Wqx08nztZ1omqSRmQfWt8dimYndyJ5f4jUkAIfzKvyPtgvpXJo4dvWVpAeNUiPoLR3dcFjrMz9b_EkfXFI3NS0hkSqwFytazHa3v_o_V7TpIi5XBJRgfba5pYJlDeqnEPhEEScOD_jhSTjBDVJWp2j6iEsxmOJ5KKCzyG4FX0laJO16lcGcBsUzkVvIIH2n5FZXyrqMcL2Evgko6d64VAtx4wA5Kxanvs6igS3bIV7MUjpeP67aTcAafVCR20hgWT_IpZ6qaBidr68H6dzyeHTMYfeE-2NvO5o8mrNg65jndIDGA84zyFkxcaqh4ghQ",
		"does not contain key id",
		0},

	jwtParcingTest{
		"NEG: non RSA256 alg, invalid sig, 1 group",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiIsIng1dCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSIsImtpZCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSJ9.eyJhdWQiOiJhcGk6Ly8wNDJkODA3ZC1hMThiLTQ5NjUtYTgyNC1lZmY3Mzg1NjA3ZTYiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC8yM2Q0MjUwZS0zNzExLTQ4ZTEtOGRmZS1kMjdmMjg1MTU0YTEvIiwiaWF0IjoxNjY1ODUyMzYxLCJuYmYiOjE2NjU4NTIzNjEsImV4cCI6MTY2NTg1NjI2MSwiYWlvIjoiRTJaZ1lJajc5dXFmdzNrMlZ4bGg4MzJMQlhZb0F3QT0iLCJhcHBpZCI6IjUzNTllNmEzLWRjMWUtNGFjMS04MDczLWQwZTY1OGEwMDJjNCIsImFwcGlkYWNyIjoiMSIsImdyb3VwcyI6WyIyZTEyZDFiYi1jNDhlLTQyY2QtYjRjZC01NzMzODVlNGE5ZGMiXSwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvMjNkNDI1MGUtMzcxMS00OGUxLThkZmUtZDI3ZjI4NTE1NGExLyIsIm9pZCI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInJoIjoiMC5BWUlBRGlYVUl4RTM0VWlOX3RKX0tGRlVvWDJBTFFTTG9XVkpxQ1R2OXpoV0ItYVZBQUEuIiwicm9sZXMiOlsiYXJnb25BZG1pblJvbGUxIiwidGVzdFJvbGUxIl0sInN1YiI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInRpZCI6IjIzZDQyNTBlLTM3MTEtNDhlMS04ZGZlLWQyN2YyODUxNTRhMSIsInV0aSI6InhSZ09ydFRZLWsydHE4M0JfTEVCQUEiLCJ2ZXIiOiIxLjAifQ.vlyvdWmVb1lnRzZg6vtAGOh9JSwWKK_6QE_bX19fmELbgL9WlmZOKpwNItqwLTS-JVuLG0qjFs7Atl1wUtXt6o5xrxOt5xOTB3YYw0ecJ6LbORd2cYqHM2vYln0mejq3soVmOjQNvaz9b17bewlw-A_UtMKLVK2WphC9KtLurazKqkS_wg3ZA3BWhFG9qZKLkbC-TSUKGboDfzZ8zcSXYjAioXkBnTqMcVngfHOcPX0VGfn9Stv9k5bonMwNQHcFNw1TND_JpaFVLMxpTn3p8JhRj_YplW70a5I2n-N5JUjW_2KVLQl122kGyxfHHSq0g55gvVSWNFfMg6U0ScLMUA",
		"unexpected signing method:",
		0},

	jwtParcingTest{
		"NEG: invalid signature",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSIsImtpZCI6IjJaUXBKM1VwYmpBWVhZR2FYRUpsOGxWMFRPSSJ9.eyJhdWQiOiJhcGk6Ly8wNDJkODA3ZC1hMThiLTQ5NjUtYTgyNC1lZmY3Mzg1NjA3ZTYiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC8yM2Q0MjUwZS0zNzExLTQ4ZTEtOGRmZS1kMjdmMjg1MTU0YTEvIiwiaWF0IjoxNjY1OTk5MzI2LCJuYmYiOjE2NjU5OTkzMjYsImV4cCI6MTY2NjAwMzIyNiwiYWlvIjoiRTJaZ1lLaStuWnA3NmQ0TDdRc0Y2Zzk1TGF3NUFBPT0iLCJhcHBpZCI6IjUzNTllNmEzLWRjMWUtNGFjMS04MDczLWQwZTY1OGEwMDJjNCIsImFwcGlkYWNyIjoiMSIsImdyb3VwcyI6WyI0YzllMzhjZS1hZjY3LTRjMDgtYmRlZS1iZThmMDZiZmFjNmIiLCIyZTEyZDFiYi1jNDhlLTQyY2QtYjRjZC01NzMzODVlNGE5ZGMiLCJjNzExM2EzNi1kZTRkLTRmY2MtOWVhMC01MjM3MjcxMGQzY2QiXSwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvMjNkNDI1MGUtMzcxMS00OGUxLThkZmUtZDI3ZjI4NTE1NGExLyIsIm9pZCI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInJoIjoiMC5BWUlBRGlYVUl4RTM0VWlOX3RKX0tGRlVvWDJBTFFTTG9XVkpxQ1R2OXpoV0ItYVZBQUEuIiwicm9sZXMiOlsiYXJnb25BZG1pblJvbGUxIiwidGVzdFJvbGUxIl0sInN1YiI6ImM2MDNhZmIwLWNjODUtNDYzMy1hYWM3LWI0YWJjODQzZjIwMyIsInRpZCI6IjIzZDQyNTBlLTM3MTEtNDhlMS04ZGZlLWQyN2YyODUxNTRhMSIsInV0aSI6Im1SNmhEdkZjWmt1YTdueGdvSkFsQUEiLCJ2ZXIiOiIxLjAifQ.XXXXXjm24js9BSGw6HODQ8V6GwVeWpsU0G9i6CFsyRcJeEPXGLiE5GnTgYpsl40AMIL_uxwRgr_IObt50uIzbdmzZ_9GO_XZLUF_6oH88-UD2M4HNZVkQL3oRluloiv41hFrUPHzTTCxJ0tUKO7YjEuPAeL_To9RQ-Cu_Jdfc5SPwR3KAN1K5numTXZ1szFbI2q0S7eDO9WmhV8XDxO-E9fFI2Bdr2fG8wGQYkJhsPny5of4rt2sREfB1RHDpEHOtDelAp2jfkbiAbP_y4NqDMI8K-7T5jp620Nqj-AUXaTgubTX-V24FWWOtwUkzM606UW-Fm2P_22JXig5Qq35vg",
		"crypto/rsa: verification error",
		0},

	jwtParcingTest{
		"NEG: token format error",
		"not_a.token.format.1232",
		"token format error",
		0},
}

func prolog(t *testing.T) {
	ok := os.Getenv("MSAD")
	if ok == "" {
		t.Skip("MS AD not available, skip")
	}

	testconfig.Set(t)
}

func TestToken(t *testing.T) {
	prolog(t)

	jwk.InitJWKCache()

	for _, tokenCase := range jwtParcingTests {

		t.Run(tokenCase.testDesc, func(t *testing.T) {

			testToken, err := token.NewToken([]byte(tokenCase.tokenStr))
			var testResult bool

			if err != nil {
				if strings.Contains(fmt.Sprint(err), tokenCase.expectedErr) && tokenCase.gropuCount == len(testToken.Groups) {
					testResult = true
				} else {
					testResult = false
					t.Errorf("Error while parsing token: %s", err)
				}
			} else {
				if tokenCase.gropuCount == len(testToken.Groups) {
					testResult = true
				} else {
					testResult = false
				}
			}

			assert.Equal(t, true, testResult, tokenCase.testDesc)

		})
	}
}