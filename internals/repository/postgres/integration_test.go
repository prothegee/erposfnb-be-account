package internal_repository_postgres_test

import (
	"context"
	"os"
	"testing"

	iC "github.com/prothegee/erposfnb-be-account/internals/containers"
	// internal repository postgres
	irPG "github.com/prothegee/erposfnb-be-account/internals/repository/postgres"

	"github.com/stretchr/testify/assert"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// --------------------------------------------------------- //

const (
	postgresContainerName        = "test-pg18-erposfnb-be-account"
	postgresReplicaContainerName = "test-pg18-erposfnb-be-account-replica"
)

// Relative to this test package directory
var postgresMigrationFiles = []string{
	"../../../migrations/000001_schema.sql",
	"../../../migrations/000002_table_account_users.sql",
	"../../../migrations/000003_table_account_information.sql",
}

// --------------------------------------------------------- //

// Test integration for postgres container
//
// One gorm instance reads and writes through dbresolver,
// the primary serves writes, the replica serves reads.
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
		irPG.PostgreSQL_PORT_PRIMARY,
		"",
		"",
		"",
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

	// 4. Primary schema, the migration sql files run on the primary
	for _, file := range postgresMigrationFiles {
		migration, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("fail to read migration %s: %v", file, err)
		}
		err = db.Exec(string(migration)).Error
		if err != nil {
			t.Fatalf("run migration %s: %v", file, err)
		}
	}

	// 5. Create replica container, the same sql files seed the replica schema
	replicaCtr, err := iC.ContainerPostgresStartWithInit(
		ctx,
		irPG.PostgreSQL_IMAGE_VERSION,
		postgresReplicaContainerName,
		irPG.PostgreSQL_PORT_REPLICA,
		"",
		"",
		"",
		postgresMigrationFiles...,
	)
	if err != nil {
		t.Fatalf("fail to start replica container: %v", err)
	}
	t.Cleanup(func() {
		_ = replicaCtr.Terminate(ctx)
	})

	// 6. Replica dsn
	replicaDSN, err := replicaCtr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("replica connection string: %v", err)
	}

	// 7. Register dbresolver, primary source for write, replica for read
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{gormpg.Open(replicaDSN)},
	}))
	if err != nil {
		t.Fatalf("register dbresolver: %v", err)
	}

	// 8. Insert user, gorm routes the create to the primary
	user := irPG.UsersRow{
		Email:    "primary-replica@example.com",
		Password: "hashed-password",
	}
	err = db.
		Omit("id", "id_ref", "username", "pub_path", "pub_pict").
		Create(&user).Error
	if err != nil {
		t.Fatalf("insert user into primary: %v", err)
	}

	// 9. Read user on the default route, the replica serves the read
	var fromReplica irPG.UsersRow
	err = db.Where("email = ?", user.Email).First(&fromReplica).Error
	if !assert.ErrorIs(t, err, gorm.ErrRecordNotFound) {
		t.Fatalf("replica must not hold the primary row: %v", err)
	}

	// 10. Read user with the write clause, the primary holds the inserted user
	var fromPrimary irPG.UsersRow
	err = db.Clauses(dbresolver.Write).Where("email = ?", user.Email).First(&fromPrimary).Error
	if err != nil {
		t.Fatalf("read user from primary: %v", err)
	}
	if !assert.NotEmpty(t, fromPrimary.ID) {
		t.Fatalf("primary row has no id")
	}

	t.Logf("ready: container=%s replica=%s", ctr.GetContainerID(), replicaCtr.GetContainerID())
}
