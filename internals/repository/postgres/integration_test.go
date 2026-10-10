package internal_repository_postgres_test

import (
	"context"
	"os"
	"testing"

	iC "github.com/prothegee/erposfnb-be-account/internals/containers"
	// internal repository postgres
	irPG "github.com/prothegee/erposfnb-be-account/internals/repository/postgres"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// --------------------------------------------------------- //

const postgresContainerName = "test-pg18-erposfnb-be-account"

// --------------------------------------------------------- //

// Test integration for postgres container
//
// Note that the init schema and table with gorm can't use complex command,
// so it required valid sql file
func Test_Postgres_Container(t *testing.T) {
	ctx := context.Background()

	// 1. Create container
	ctr, err := iC.ContainerPostgresStart(
		ctx,
		irPG.PostgreSQL_IMAGE_VERSION,
		postgresContainerName,
	)
	if err != nil {
		t.Fatalf("fail to start container: %v", err)
	}
	t.Cleanup(func() {
		_ = ctr.Terminate(ctx)
	})

	// 2. Container dsn
	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	// 3. GORM
	db, err := gorm.Open(gormpg.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}

	// 4. Create schema
	schema := "../../../migrations/000001_schema.sql"
	schemaInit, err := os.ReadFile(schema)
	if err != nil {
		t.Fatalf("fail to read schema sql file: %v", err)
	}
	err = db.Exec(string(schemaInit)).Error
	if err != nil {
		t.Fatalf("create schema: %v", err)
	}

	// 5. Create table account.users
	tableUsers := "../../../migrations/000002_table_account_users.sql"
	tableUsersInit, err := os.ReadFile(tableUsers)
	if err != nil {
		t.Fatalf("fail to read schema sql file: %v", err)
	}
	err = db.Exec(string(tableUsersInit)).Error
	if err != nil {
		t.Fatalf("create schema: %v", err)
	}

	// 6. Create table account.information
	tableInformation := "../../../migrations/000003_table_account_information.sql"
	tableInformationInit, err := os.ReadFile(tableInformation)
	if err != nil {
		t.Fatalf("fail to read schema sql file: %v", err)
	}
	err = db.Exec(string(tableInformationInit)).Error
	if err != nil {
		t.Fatalf("create schema: %v", err)
	}

	t.Logf("ready: container=%s", ctr.GetContainerID())
}
