package config_test

import (
	"testing"

	"admincheckapi/api/config"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("config empty string", func(t *testing.T) {
		var input []byte
		_, err := config.NewSetupValueSet(input)
		if err == nil {
			t.Fatalf("No error loading empty string as setup")
		}
	})

	t.Run("config invalid string", func(t *testing.T) {
		var input []byte = []byte("whatever")
		_, err := config.NewSetupValueSet(input)
		if err == nil {
			t.Fatalf("No error loading invalid string as setup")
		}
	})

	t.Run("config invalid yaml string", func(t *testing.T) {
		var input []byte = []byte(
			`providers:
msad:
kind: msad
env:
tenant_id: abc
client_id: 3b0c33ac-93c9-4f6e-bc70-d508465155a0
servers:
- http:
  kind: http
  env:
    port: 11
    address: xxx
- http:
  kind: https
  env:
    port: 22
    address: yyy
backends:
postgres:
  kind: postgres
  env:
    user: test
    pass: test
    dbname: argonadmindb
    host: localhost`)
		_, err := config.NewSetupValueSet(input)
		if err == nil {
			t.Fatalf("Error loading valid setup")
		}
	})

	t.Run("config valid string", func(t *testing.T) {
		var input []byte = []byte(
			`providers:
- msad:
  kind: msad
  env:
    tenant_id: abc
    client_id: 3b0c33ac-93c9-4f6e-bc70-d508465155a0
servers:
- http:
  kind: http
  env:
    port: 11
    address: xxx
- http:
  kind: https
  env:
    port: 22
    address: yy
backends:
- postgres:
  kind: postgres
  env:
    user: test
    pass: test
    dbname: argonadmindb
    host: localhost`)
		s, err := config.NewSetupValueSet(input)
		if err != nil {
			t.Fatalf("Error loading valid setup")
		}
		assert.Equal(t, s.UsedBackend, "postgres")
		assert.Equal(t, s.TenantId, "abc")
		assert.Equal(t, s.ServerPort, "11")
		assert.Equal(t, s.ServerIPAddress, "xxx")
	})
}
