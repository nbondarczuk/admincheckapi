package azure

//
// BackendCredentialsAzure stores only the JWT token
//
type BackendCredentialsAzure struct {
	token string
}

//
// NewBackendCredentials
//
func NewBackendCredentials(t string) (BackendCredentialsAzure, error) {
	return BackendCredentialsAzure{token: t}, nil
}

//
// ConnectString uses token to implement connection secret
//
func (bc BackendCredentialsAzure) ConnectString() string {
	return bc.token
}

