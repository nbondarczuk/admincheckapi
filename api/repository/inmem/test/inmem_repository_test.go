package repository_test

import (
	"testing"

	"admincheckapi/api/repository"

	"github.com/stretchr/testify/assert"
)

func TestInMemNewClientRepository(t *testing.T) {
	// Invalid repository kind causes error
	t.Run("invalid repository kind requested", func(t *testing.T) {
		_, err := repository.NewClientAdminGroupRepository("whatever")
		if err == nil {
			t.Fatalf("Error creating invalid repository")
		}
	})

	// Created inmem repository closes
	t.Run("closes created repository", func(t *testing.T) {
		r, err := repository.NewClientAdminGroupRepository("inmem")
		if err != nil {
			t.Fatalf("Error creating repository: %s", err.Error())
		}

		r.Close()
	})

	t.Run("create one group", func(t *testing.T) {
		r, err := repository.NewClientAdminGroupRepository("inmem")
		if err != nil {
			t.Fatalf("Error creating repository: %s", err.Error())
		}

		cags, count, err := r.CreateClientGroup("client", "group")
		if err != nil {
			t.Fatalf("Error creating client group: %s", err.Error())
		}
		if count != 1 {
			t.Fatalf("Invalid returned created client groups count: %d", count)
		}
		if len(cags) != 1 {
			t.Fatalf("Invalid returned created client group slice length: %d", len(cags))
		}

		r.PurgeClientGroups()
		r.Close()
	})

	t.Run("count created one group", func(t *testing.T) {
		r, err := repository.NewClientAdminGroupRepository("inmem")
		if err != nil {
			t.Fatalf("Error creating repository: %s", err.Error())
		}

		cags, count, err := r.CreateClientGroup("client", "group")
		if err != nil {
			t.Fatalf("Error creating client group: %s", err.Error())
		}
		if count != 1 {
			t.Fatalf("Invalid returned created client groups count: %d", count)
		}
		if len(cags) != 1 {
			t.Fatalf("Invalid returned created client group slice length: %d", len(cags))
		}

		count, err = r.CountClientGroups("client", "group")
		if err != nil {
			t.Fatalf("Error creating client group: %s", err.Error())
		}
		if count != 1 {
			t.Fatalf("Invalid numer of counter groups: %d", count)
		}

		r.PurgeClientGroups()
		r.Close()
	})

	// Crearte inmem repository, create 1 group and read it
	t.Run("created one group read returns same group", func(t *testing.T) {
		r, err := repository.NewClientAdminGroupRepository("inmem")
		if err != nil {
			t.Fatalf("Error creating repository: %s", err.Error())
		}

		cags, count, err := r.CreateClientGroup("client", "group")
		if err != nil {
			t.Fatalf("Error creating client group: %s", err.Error())
		}
		if count != 1 {
			t.Fatalf("Invalid returned created client groups count: %d", count)
		}
		if len(cags) != 1 {
			t.Fatalf("Invalid returned created client group slice length: %d", len(cags))
		}
		assert.Equal(t, cags[0].Client, "client")
		assert.Equal(t, cags[0].AdminGroupId, "group")

		cags, count, err = r.ReadClientGroups("client")
		if err != nil {
			t.Fatalf("Error reading client groups: %s", err.Error())
		}
		if count != 1 {
			t.Fatalf("Invalid returned read client groups count: %d", count)
		}
		if len(cags) != 1 {
			t.Fatalf("Invalid returned read client group slice length: %d", len(cags))
		}
		assert.Equal(t, cags[0].Client, "client")
		assert.Equal(t, cags[0].AdminGroupId, "group")

		r.PurgeClientGroups()
		r.Close()
	})

	t.Run("delete created one group read returns no groups", func(t *testing.T) {
		r, err := repository.NewClientAdminGroupRepository("inmem")
		if err != nil {
			t.Fatalf("Error creating repository: %s", err.Error())
		}

		cags, count, err := r.CreateClientGroup("client", "group")
		if err != nil {
			t.Fatalf("Error creating client group: %s", err.Error())
		}
		if count != 1 {
			t.Fatalf("Invalid returned created client groups count: %d", count)
		}
		if len(cags) != 1 {
			t.Fatalf("Invalid returned created client group slice length: %d", len(cags))
		}
		assert.Equal(t, cags[0].Client, "client")
		assert.Equal(t, cags[0].AdminGroupId, "group")

		cags, count, err = r.DeleteClientGroup("client", "group")
		if err != nil {
			t.Fatalf("Error deleting client group: %s", err.Error())
		}
		if count != 1 {
			t.Fatalf("Invalid returned deleted client groups count: %d", count)
		}

		cags, count, err = r.ReadClientGroups("client")
		if err != nil {
			t.Fatalf("Error reading client group: %s", err.Error())
		}
		if count != 0 {
			t.Fatalf("Invalid returned read client groups count: %d", count)
		}
		if len(cags) != 0 {
			t.Fatalf("Invalid returned read client group slice length: %d", len(cags))
		}

		r.PurgeClientGroups()
		r.Close()
	})

}
