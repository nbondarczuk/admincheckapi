package gorm_repository_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	bp "admincheckapi/api/backend/postgres"
	"admincheckapi/api/model"
	rg "admincheckapi/api/repository/gorm"
	"admincheckapi/test/testconfig"
)

func TestNewClientAdminGroupRepository(t *testing.T) {
	testconfig.Set(t)

	var mock sqlmock.Sqlmock
	var r rg.GORMClientRepository
	t.Run("first mock gorm", func(t *testing.T) {
		var (
			mocksqldb *sql.DB
			err       error
		)

		mocksqldb, mock, err = sqlmock.New()
		if err != nil {
			t.Errorf("Failed to open mock sql db, got error: %v", err)
		}
		if mock == nil {
			t.Errorf("Failed to make mock, nil")
		}

		mockbackend := bp.BackendPostgres{
			Kind:          "postgres",
			ConnectString: "mock",
			Sqldb:         mocksqldb,
		}

		mockdialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 mocksqldb,
			PreferSimpleProtocol: true,
		})

		r, err = rg.NewClientAdminGroupRepository(mockbackend, mockdialector)
		if err != nil {
			t.Fatalf("Error creating gorm repository: %s", err)
		}
	})

	t.Run("read client groups", func(t *testing.T) {
		const (
			size   = int64(10)
			client = "client"
		)
		var now = time.Now()

		models := make([]model.ClientAdminGroup, size)
		for i := int64(0); i < size; i++ {
			models[i] = model.ClientAdminGroup{
				Client:     client,
				AdminGroupId: fmt.Sprintf("group_%d", i),
				Model: gorm.Model{
					ID:        uint(i),
					CreatedAt: now,
					UpdatedAt: now,
				},
			}
		}

		rows := sqlmock.NewRows([]string{"client", "AdminGroupId", "id", "CreatedAt", "UpdatedAt", "DeletedAt"})
		for _, model := range models {
			rows.AddRow(model.Client, model.AdminGroupId, model.ID, model.CreatedAt, model.UpdatedAt, nil)
		}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "client_admin_groups" WHERE client = $1 AND "client_admin_groups"."deleted_at" IS NULL`)).
			WithArgs(client).
			WillReturnRows(rows)
		ret, i, err := r.ReadClientGroups(client)
		assert.NoError(t, err)
		assert.Equal(t, size, i)
		assert.Equal(t, models, ret)
	})

	t.Run("count client groups", func(t *testing.T) {
		const (
			client = "client"
			group  = "group"
			size   = int64(10)
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "client_admin_groups" WHERE client = $1 AND admin_group_id = $2 AND "client_admin_groups"."deleted_at" IS NULL`)).
			WithArgs(client, group).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(size))
		ret, err := r.CountClientGroups(client, group)
		assert.NoError(t, err)
		assert.Equal(t, size, ret)
	})

	r.Close()
}
