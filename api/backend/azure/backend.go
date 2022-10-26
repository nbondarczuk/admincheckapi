package azure

//
// BackendAzure must keep JWT token for API requests
//
type BackendAzure struct {
	Kind, Token string
}

//
// NewBackendAzure stores thetoken for subsequent usage
//
func NewBackend(token string) (BackendAzure, error) {
	return BackendAzure{
			Kind:  "azure",
			Token: token,
		},
		nil
}

//
// Version provides Azure API version
//
func (be BackendAzure) Version() (string, error) {
	return "1.0", nil
}

//
// Ping checks the trivial API method
//
func (b BackendAzure) Ping() error {
	return nil
}

//
// Credentials
//
func (b BackendAzure) Credentials() string {
	return b.Token
}

//
// Close
//
func (b BackendAzure) Close() {
	b.Token = ""
}
