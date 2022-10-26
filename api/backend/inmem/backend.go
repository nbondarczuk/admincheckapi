package inmem

type BackendInMem struct {}

type BackendCredentialsInMem struct{}

func (bc BackendCredentialsInMem) ConnectString() string {
	return ""
}

func NewBackend() (BackendInMem, error) {
	return BackendInMem{}, nil
}

func (b BackendInMem) Version() (string, error) {
	return "n/a", nil
}

func (b BackendInMem) Ping() (err error) { return }

func (b BackendInMem) Credentials() string {
	return ""
}

func (b BackendInMem) Close() {}
