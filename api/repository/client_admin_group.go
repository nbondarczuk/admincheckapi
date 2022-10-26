package repository

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	
	backend "admincheckapi/api/backend"
	backendmysql "admincheckapi/api/backend/mysql"
	backendpostgres "admincheckapi/api/backend/postgres"
	"admincheckapi/api/model"
	"admincheckapi/api/repository/gorm"
	"admincheckapi/api/repository/inmem"
)

// ClientAdminGroupRepository
type ClientAdminGroupRepository interface {
	CountClientGroups(client, group string) (int64, error)
	ReadClientGroups(client string) ([]model.ClientAdminGroup, int64, error)
	CreateClientGroup(client, group string) ([]model.ClientAdminGroup, int64, error)
	CreateClientGroups(client string, groups []model.ClientAdminGroup) ([]model.ClientAdminGroup, int64, error)	
	DeleteClientGroup(client, group string) ([]model.ClientAdminGroup, int64, error)
	PurgeClientGroups() error
	Close()
}

//
// NewClientRepository is a main dispatch, inmem or something else (gorm backends, really)
// It hides the type of backend db which is a detail. A new backend of another
// GORM dblike mYSQLcan be passed as parameter and GORM will handle it.
//
func NewClientAdminGroupRepository(kind string) (ClientAdminGroupRepository, error) {
	b, err := backend.NewBackend(kind)
	if err != nil {
		return nil, fmt.Errorf("Error creating backend: %s", err)
	}

	// dispatch for repository kind
	if kind == "mysql" {
		return gorm.NewClientAdminGroupRepository(b,
			mysql.New(mysql.Config{Conn: b.(backendmysql.BackendMySQL).Sqldb}))
	} else if kind == "postgres" {
		return gorm.NewClientAdminGroupRepository(b,
			postgres.New(postgres.Config{Conn: b.(backendpostgres.BackendPostgres).Sqldb}))
	} else if kind == "inmem" {
		return inmem.NewClientAdminGroupRepository(b)
	}

	return nil, fmt.Errorf("Invalid kind of repository: %s", kind)
}
