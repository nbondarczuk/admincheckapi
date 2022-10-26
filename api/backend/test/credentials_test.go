package backend_test

import (
	"os"
	"testing"

	"admincheckapi/api/backend"

	"github.com/stretchr/testify/assert"
)

func initPostgresEnv() {
	os.Setenv("POSTGRES_USER", "test")
	os.Setenv("POSTGRES_PASS", "test")
	os.Setenv("POSTGRES_DBNAME", "testdb")
	os.Setenv("POSTGRES_HOST", "localhost")
}

func resetPostgresEnv() {
	os.Unsetenv("POSTGRES_USER")
	os.Unsetenv("POSTGRES_PASS")
	os.Unsetenv("POSTGRES_DBNAME")
	os.Unsetenv("POSTGRES_HOST")
}

func TestNewBackendCredentials(t *testing.T) {
	initPostgresEnv()

	t.Run("invalid backend name causes error", func(t *testing.T) {
		_, err := backend.NewBackendCredentials("whatever")
		if err == nil {
			t.Fatalf("No error creating invalid credentials: %s", err.Error())
		}
	})

	t.Run("postgres credentials created", func(t *testing.T) {
		bc, err := backend.NewBackendCredentials("postgres")
		if err != nil {
			t.Fatalf("Error creating postgres credentials: %s", err.Error())
		}
		cs := bc.ConnectString()
		assert.Equal(t,
			"host=localhost port=5432 user=test password=test dbname=testdb sslmode=disable",
			cs)
	})

	resetPostgresEnv()
}
