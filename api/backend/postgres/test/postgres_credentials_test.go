package postgres_test

import (
	"os"
	"testing"

	"admincheckapi/api/backend/postgres"
	"admincheckapi/test/testconfig"

	"github.com/stretchr/testify/assert"
)

func prolog(t *testing.T) {
	ok := os.Getenv("POSTGRES")
	if ok == "" {
		t.Skip("Postgres DB not available, skip")
	}

	testconfig.Set(t)
}

func TestNewBackendCredentialsPostgres(t *testing.T) {
	prolog(t)

	t.Run("postgres credentials created", func(t *testing.T) {
		bc, err := postgres.NewBackendCredentials()
		if err != nil {
			t.Fatalf("Error creating postgres credentials: %s", err.Error())
		}
		cs := bc.ConnectString()
		assert.Equal(t,
			cs,
			"host=localhost port=5432 user=test password=test dbname=argonadmindb sslmode=disable")
	})
}
