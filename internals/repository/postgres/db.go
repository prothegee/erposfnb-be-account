package internal_repository_postgres

const (
	// postgres:18
	PostgreSQL_IMAGE_VERSION = "postgres:18"
	// Port primary
	PostgreSQL_PORT_PRIMARY = 5432
	// Port replica
	PostgreSQL_PORT_REPLICA = 5433
	// Default db name both for primary and replica
	PostgreSQL_DB_NAME = "erposfnb_be_account"
)
