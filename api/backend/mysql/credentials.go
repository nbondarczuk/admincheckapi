package mysql

import (
	"fmt"
	"os"
)

const (
	DefaultMySQLPort string = "3306"
)

//BackendCredentialsMySQL is a standard set of required login credentials
type BackendCredentialsMySQL struct {
	user, password, dbname, host, port string
}

//
// NewBackendCredentialsMySQL build an interfane respresentation of a connect string
//
func NewBackendCredentials() (BackendCredentialsMySQL, error) {
	user := os.Getenv("MYSQL_USER")
	if user == "" {
		return BackendCredentialsMySQL{},
			fmt.Errorf("Missing env variable: %s", "MYSQL_USER")
	}

	password := os.Getenv("MYSQL_PASS")
	if password == "" {
		return BackendCredentialsMySQL{},
			fmt.Errorf("Missing env variable: %s", "MYSQL_PASS")
	}

	dbname := os.Getenv("MYSQL_DBNAME")
	if dbname == "" {
		return BackendCredentialsMySQL{},
			fmt.Errorf("Missing env variable: %s", "MYSQL_DBNAME")
	}

	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		return BackendCredentialsMySQL{},
			fmt.Errorf("Missing env variable: %s", "MYSQL_HOST")
	}

	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = DefaultMySQLPort
	}

	return BackendCredentialsMySQL{
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
// to be use in the DB connection like: user:pass@tcp(127.0.0.1:3306)/dbname
//
func (bc BackendCredentialsMySQL) ConnectString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		bc.user, bc.password, bc.host, bc.port, bc.dbname)
}
