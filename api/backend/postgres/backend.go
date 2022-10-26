package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"admincheckapi/api/config"
)

// Backend for Postgres DB
type BackendPostgres struct {
	Kind          string
	ConnectString string
	Sqldb         *sql.DB
}

//
// NewBackend creates and opens new Postgres DB connection with GORM layer
//
func NewBackend(bc BackendCredentialsPostgres) (BackendPostgres, error) {
	log.Trace("Begin: postgres.NewBackend")

	cs := bc.ConnectString()
	log.Debugf("Postgres DB connect string: %s", cs)

	sqldb, err := sql.Open("postgres", cs)
	if err != nil {
		return BackendPostgres{}, fmt.Errorf("Error opepning Postgres DB connection: %s", err)
	}
	log.Debug("Connected Postgres DB")

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqldb.SetMaxIdleConns(config.Setup.SQLMaxIdleConns)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqldb.SetMaxOpenConns(config.Setup.SQLMaxOpenConns)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqldb.SetConnMaxLifetime(config.Setup.SQLMaxLifetime)

	log.Trace("End: postgres.NewBackend")
	return BackendPostgres{
		Kind:          "postgres",
		ConnectString: cs,
		Sqldb:         sqldb,
	}, nil
}

//
// Version obtains the backend server version: it is highly database dependent
//
func (b BackendPostgres) Version() (string, error) {
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
func (b BackendPostgres) Ping() error {
	log.Trace("Begin: Ping")

	err := b.Sqldb.Ping()
	if err != nil {
		return fmt.Errorf("Error pinging postgres: %s", err)
	}
	log.Debug("Pinged Postgres DB")

	log.Tracef("End: Ping")
	return nil
}

//
// Credentials
//
func (b BackendPostgres) Credentials() string {
	return b.ConnectString
}

//
// Close backend connection
//
func (b BackendPostgres) Close() {
	log.Trace("Begin: Close")

	b.Sqldb.Close()
	log.Debug("Closed connection to Postgres DB")

	log.Trace("End: Close")
}
