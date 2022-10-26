package postgres

import (
	"fmt"
	"os"
)

const (
	DefaultPostgresPort string = "5432"
)

//BackendCredentialsPostgres is a standard set of required login credentials
type BackendCredentialsPostgres struct {
	user, password, dbname, host, port string
}

//
// NewBackendCredentialsPostgres build an interfane respresentation of a connect string
//
func NewBackendCredentials() (BackendCredentialsPostgres, error) {
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		return BackendCredentialsPostgres{},
			fmt.Errorf("Missing env variable: %s", "POSTGRES_USER")
	}

	password := os.Getenv("POSTGRES_PASS")
	if password == "" {
		return BackendCredentialsPostgres{},
			fmt.Errorf("Missing env variable: %s", "POSTGRES_PASS")
	}

	dbname := os.Getenv("POSTGRES_DBNAME")
	if dbname == "" {
		return BackendCredentialsPostgres{},
			fmt.Errorf("Missing env variable: %s", "POSTGRES_DBNAME")
	}

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		return BackendCredentialsPostgres{},
			fmt.Errorf("Missing env variable: %s", "POSTGRES_HOST")
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = DefaultPostgresPort
	}

	return BackendCredentialsPostgres{
			user:     user,
			password: password,
			dbname:   dbname,
			host:     host,
			port:     port,
		},
		nil
}

//
// ConnectString produces the external respresentation of the connect string
// to be use in the DB connection
//
func (bc BackendCredentialsPostgres) ConnectString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		bc.host, bc.port, bc.user, bc.password, bc.dbname)
}
