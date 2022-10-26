package backend

import (
	"fmt"

	"admincheckapi/api/backend/azure"
	"admincheckapi/api/backend/postgres"	
)

//
// BackendCredentials is an interface handling db specific connect string
//
type BackendCredentials interface {
	ConnectString() string
}

//
// NewBackendCredentials fills the structure with items required for login
//
func NewBackendCredentials(kind string) (BackendCredentials, error) {
	if kind == "inmem" {
		return nil, nil
	} else if kind == "postgres" {
		return postgres.NewBackendCredentials()
	} else if kind[:5] == "azure" {
		return azure.NewBackendCredentials(kind[6:])
	}
		
	return nil, fmt.Errorf("Invalid kind of backend: " + kind)
}
