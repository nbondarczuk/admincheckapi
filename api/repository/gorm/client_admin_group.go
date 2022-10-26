package gorm

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"admincheckapi/api/backend"
	"admincheckapi/api/config"
	"admincheckapi/api/model"

	log "github.com/sirupsen/logrus"
)

// GORM Client handle
type GORMClientRepository struct {
	be     backend.Backend
	gormdb *gorm.DB
}

//
// NewClientRepository creates a handle for domain operations on a client using gorm
//
func NewClientAdminGroupRepository(b backend.Backend, dial gorm.Dialector) (GORMClientRepository, error) {
	log.Trace("Begin: NewClientAdminGroupRepository")

	log.Debug("Pinging backend DB")
	err := b.Ping()
	if err != nil {
		return GORMClientRepository{},
			fmt.Errorf("Error pinging  backend DB: %s", err)
	}
	log.Debug("Pinged backend DB")

	var c gorm.Config
	if config.Setup.LogGORM {
		c = gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	} else {
		c = gorm.Config{
			Logger: GormLogger{LogLevel: logger.Silent},
		}
	}

	log.Debug("Opening GORM on backend DB")
	gormdb, err := gorm.Open(dial, &c)
	if err != nil {
		return GORMClientRepository{},
			fmt.Errorf("Error opening GORM on backend DB: %s", err)
	}
	log.Debug("Opened GORM on ackend DB")

	log.Debug("Migrating schema to GORM")
	gormdb.AutoMigrate(&model.ClientAdminGroup{})
	log.Debug("Migrated schema to GORM")

	log.Trace("End: NewClientAdminGroupRepository")
	return GORMClientRepository{b, gormdb}, nil
}

//
// ReadClientGroups reads all groups of the client
//
func (r GORMClientRepository) CountClientGroups(client, group string) (int64, error) {
	log.Trace("Begin: CountClientGroups")
	var count int64
	result := r.gormdb.Model(&model.ClientAdminGroup{}).
		Where("client = ?", client).
		Where("admin_group_id = ?", group).
		Count(&count)
	log.Trace("End: CountClientGroups")
	return count, result.Error
}

//
// ReadClientGroups reads all groups of the client
//
func (r GORMClientRepository) ReadClientGroups(client string) ([]model.ClientAdminGroup, int64, error) {
	log.Trace("Begin: ReadClientGroups")
	var cgs []model.ClientAdminGroup
	result := r.gormdb.Find(&cgs, "client = ?", client)
	log.Trace("End: ReadClientGroups")
	return cgs, result.RowsAffected, result.Error
}

//
// CreateClientGroup creates mappig between client and a group
//
func (r GORMClientRepository) CreateClientGroup(client, group string) ([]model.ClientAdminGroup, int64, error) {
	log.Trace("Begin: CreateClientGroup")
	result := r.gormdb.Create(&model.ClientAdminGroup{Client: client, AdminGroupId: group})
	log.Trace("End: CreateClientGroup")
	return []model.ClientAdminGroup{model.ClientAdminGroup{Client: client, AdminGroupId: group}},
		result.RowsAffected,
		result.Error
}

//
// CreateClientGroups creates mappig between client and many groups
//
func (r GORMClientRepository) CreateClientGroups(client string, groups []model.ClientAdminGroup) ([]model.ClientAdminGroup, int64, error) {
	log.Trace("Begin: CreateClientGroups")
	result := r.gormdb.Create(&groups)
	log.Trace("End: CreateClientGroups")
	return groups,
		result.RowsAffected,
		result.Error
}

//
// DeleteClientGroup deletes mappng between client and a group
//
func (r GORMClientRepository) DeleteClientGroup(client, group string) ([]model.ClientAdminGroup, int64, error) {
	log.Trace("Begin: DeleteClientGroup")
	result := r.gormdb.
		Where("client = ?", client).
		Where("admin_group_id = ?", group).
		Delete(&model.ClientAdminGroup{})
	log.Trace("End: DeleteClientGroup")
	return []model.ClientAdminGroup{},
		result.RowsAffected,
		result.Error
}

//
// PurgeClientGroups is a test only utility
//
func (r GORMClientRepository) PurgeClientGroups() error {
	log.Trace("Begin: PurgeClientGroups")
	result := r.gormdb.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Unscoped().
		Delete(&model.ClientAdminGroup{})
	log.Trace("End: PurgeClientGroups")
	return result.Error
}

//
// Close
//
func (r GORMClientRepository) Close() {
	log.Trace("Begin: Close")
	r.be.Close()
	log.Trace("End: Close")
}
