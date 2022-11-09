package azure

import (
	"admincheckapi/api/backend"
	"admincheckapi/api/graph"
)

type AzureClientRepository struct {
	token  string
	caller graph.Caller
}

//
// NewAzureClientRepository creates a handle for domain operations on a client
//
func NewClientAdminGroupRepository(b backend.Backend) (AzureClientRepository, error) {
	return AzureClientRepository{
		token: b.Credentials(),
		caller: graph.Caller{
			Token: b.Credentials(),
			URL:   graph.MSGraphURL,
		},
	}, nil
}

//
// ClientGroupName
//
func (r AzureClientRepository) ClientGroupName(id string) (string, error) {
	return r.caller.GroupName(id)
}

//
// ClientGroupId
//
func (r AzureClientRepository) ClientGroupId(name string) (string, error) {
	return r.caller.GroupId(name)
}

