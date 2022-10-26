package backend

import (
	"fmt"

	"admincheckapi/api/backend/azure"
	"admincheckapi/api/backend/inmem"
	"admincheckapi/api/backend/mysql"
	"admincheckapi/api/backend/postgres"
)

//
// Backend is an interface providing access to specific kind of db if needed
//
type Backend interface {
	Version() (string, error)
	Ping() error
	Credentials() string
	Close()
}

//
// NewBackend is a factory producing specific kind of backend db handlers based on dispatch
//
func NewBackend(kind string) (Backend, error) {
	bc, err := NewBackendCredentials(kind)
	if err != nil {
		return nil, fmt.Errorf("Can't get backend credentials: " + err.Error())
	}

	if kind == "mysql" {
		return mysql.NewBackend(bc.(mysql.BackendCredentialsMySQL))
	} else if kind == "postgres" {
		return postgres.NewBackend(bc.(postgres.BackendCredentialsPostgres))
	} else if kind[:5] == "azure" {
		return azure.NewBackend(bc.ConnectString())
	} else if kind == "inmem" {
		return inmem.NewBackend()
	}

	return nil, fmt.Errorf("Invalid kind of backend: " + kind)
}
