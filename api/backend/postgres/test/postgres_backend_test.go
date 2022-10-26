package postgres_test

import (
	"testing"

	"admincheckapi/api/backend/postgres"

	"github.com/stretchr/testify/assert"
)

func TestPostgresBackend(t *testing.T) {
	prolog(t)

	t.Run("create, check, ping, get version, close postgres db connection", func(t *testing.T) {
		cs, err := postgres.NewBackendCredentials()
		if err != nil {
			t.Fatalf("Error creating postgres connect string: %s", err.Error())
		}

		db, err := postgres.NewBackend(cs)
		if err != nil {
			t.Fatalf("Error creating repository: %s", err.Error())
		}

		assert.Equal(t, db.Kind, "postgres")

		err = db.Ping()
		if err != nil {
			t.Fatalf("Error pinging backend: %s", err.Error())
		}

		_, err = db.Version()
		if err != nil {
			t.Fatalf("Error getting backend version: %s", err.Error())
		}

		db.Close()
	})
}
