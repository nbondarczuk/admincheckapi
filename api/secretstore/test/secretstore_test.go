package secretstore_test

import (
	"os"
	"testing"

	"admincheckapi/api/secretstore"
	"admincheckapi/test/testconfig"
	//"github.com/stretchr/testify/assert"
)

func prolog(t *testing.T) {
	ok := os.Getenv("MSAD")
	if ok == "" {
		t.Skip("MS AD not available, skip")
	}

	testconfig.Set(t)
}

func TestTokenAcquire(t *testing.T) {
	prolog(t)

	t.Run("check msad token acquire", func(t *testing.T) {
		token, err := secretstore.TenantJWTToken()
		if err != nil {
			t.Errorf("Error getting token: %s", err)
		}
		if len(token) == 0 {
			t.Errorf("Invalid token length")
		}
	})
}
