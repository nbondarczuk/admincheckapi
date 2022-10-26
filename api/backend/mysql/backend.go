package mysql

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
	
	"admincheckapi/api/config"	
)

// Backend for MySQL DB
type BackendMySQL struct {
	Kind          string
	ConnectString string
	Sqldb         *sql.DB
}

//
// NewBackend creates and opens new MySQL DB connection with GORM layer
//
func NewBackend(bc BackendCredentialsMySQL) (BackendMySQL, error) {
	log.Trace("Begin: postgres.NewBackend")

	cs := bc.ConnectString()
	log.Debugf("MySQL DB connect string: %s", cs)

	sqldb, err := sql.Open("mysql", cs)
	if err != nil {
		return BackendMySQL{},
			fmt.Errorf("Error opepning MySQL DB connection: %s", err)
	}
	log.Debug("Connected MySQL DB")

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqldb.SetMaxIdleConns(config.Setup.SQLMaxIdleConns)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqldb.SetMaxOpenConns(config.Setup.SQLMaxOpenConns)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqldb.SetConnMaxLifetime(config.Setup.SQLMaxLifetime)
	
	log.Trace("End: postgres.NewBackend")
	return BackendMySQL{"mysql", cs, sqldb}, nil
}

//
// Version obtains the backend server version: it is highly database dependent
//
func (b BackendMySQL) Version() (string, error) {
	log.Trace("Begin: Version")

	var version string
	err := b.Sqldb.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		return "", fmt.Errorf("Error selection version: %s", err)
	}
	log.Debugf("Got version: %s", version)

	log.Trace("End: Version")
	return version, nil
}

//
// Close backend connection
//
func (b BackendMySQL) Ping() error {
	log.Trace("Begin: Ping")

	err := b.Sqldb.Ping()
	if err != nil {
		return fmt.Errorf("Error pinging postgres: %s", err)
	}
	log.Debug("Pinged MySQL DB")

	log.Tracef("End: Ping")
	return nil
}

//
// Credentials
//
func (b BackendMySQL) Credentials() string {
	return b.ConnectString
}

//
// Close backend connection
//
func (b BackendMySQL) Close() {
	log.Trace("Begin: Close")

	b.Sqldb.Close()
	log.Debug("Closed connection to MySQL DB")

	log.Trace("End: Close")
}
