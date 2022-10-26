package inmem

import (
	"admincheckapi/api/backend"
	"admincheckapi/api/model"
	"fmt"
)

// InMem Client handle
type InMemClientRepository struct {
	be backend.Backend
}

// Most simple implementation of in memory db: client -> (group -> int)
var db map[string][]string = make(map[string][]string)

//
// NewInMemClientRepository creates a handle for domain operations on a client using gorm
//
func NewClientAdminGroupRepository(be backend.Backend) (InMemClientRepository, error) {
	err := be.Ping()
	if err != nil {
		return InMemClientRepository{}, err
	}

	return InMemClientRepository{be}, nil
}

//
// ReadClientGroups counts the groups of the client
//
func (r InMemClientRepository) CountClientGroups(client, group string) (count int64, err error) {
	if groups, found := db[client]; found {
		for _, val := range groups {
			if val == group {
				count++
			}
		}
	}

	return
}

//
// ReadClientGroups reads all groups of the client
//
func (r InMemClientRepository) ReadClientGroups(client string) (cgs []model.ClientAdminGroup, count int64, err error) {
	cgs = make([]model.ClientAdminGroup, 0)
	count = int64(len(db[client]))
	if count > 0 {
		for _, group := range db[client] {
			cgs = append(cgs, model.ClientAdminGroup{Client: client, AdminGroupId: group})
		}
	}

	return
}

//
// CreateClientGroup creates mappig between client and a group
//
func (r InMemClientRepository) CreateClientGroup(client, group string) (cgs []model.ClientAdminGroup, count int64, err error) {
	if _, found := db[client]; !found {
		db[client] = make([]string, 0)
	}
	db[client] = append(db[client], group)

	cgs = make([]model.ClientAdminGroup, 0)
	cgs = append(cgs, model.ClientAdminGroup{Client: client, AdminGroupId: group})
	count = 1

	return
}

//
// CreateClientGroups creates mappig between client and a group
//
func (r InMemClientRepository) CreateClientGroups(client string, groups []model.ClientAdminGroup) (cgs []model.ClientAdminGroup, count int64, err error) {
	if _, found := db[client]; !found {
		db[client] = make([]string, 0)
	}
	
	for _, group := range groups {
		db[client] = append(db[client], group.AdminGroupId)
	}
	
	cgs = make([]model.ClientAdminGroup, 0)
	cgs = append(cgs, groups...)
	count = 1

	return
}

//
// DeleteClientGroup deletes mappng between client and a group
//
func (r InMemClientRepository) DeleteClientGroup(client, group string) (cgs []model.ClientAdminGroup, count int64, err error) {
	cgs = make([]model.ClientAdminGroup, 0)
	if groups, found := db[client]; found {
		i := find(groups, group)
		if i != -1 {
			db[client] = append(groups[:i], groups[i+1:]...)
			cgs = append(cgs, model.ClientAdminGroup{Client: client, AdminGroupId: group})
			count = 1
		} else {
			err = fmt.Errorf("Missing group: %s", group)
		}
	} else {
		err = fmt.Errorf("Missing client: %s", client)
	}

	return
}

//
// PurgeClientGroups is a test only function
//
func (r InMemClientRepository) PurgeClientGroups() (err error) {
	db = make(map[string][]string)
	return
}

//
// Close releases allocated resources of the repository
//
func (r InMemClientRepository) Close() {
	r.be.Close()
}

//
// find
//
func find(items []string, item string) int {
	for index, it := range items {
		if it == item {
			return index
		}
	}

	return -1
}
