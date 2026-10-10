package internal_repository_postgres

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/schema"
)

// --------------------------------------------------------- //

// Test users row column mapping
func Test_UsersRow_Mapping(t *testing.T) {
	parsed, err := schema.Parse(&UsersRow{}, &sync.Map{}, schema.NamingStrategy{})
	require.NoError(t, err, "parse users row schema")

	assert.Equal(t, "account.users", UsersRow{}.TableName())

	assert.Equal(t, "id_ref", parsed.LookUpField("IDRef").DBName)
	assert.Equal(t, "email", parsed.LookUpField("Email").DBName)
	assert.Equal(t, "phone_number", parsed.LookUpField("PhoneNumber").DBName)
	assert.Equal(t, "username", parsed.LookUpField("Username").DBName)
	assert.Equal(t, "password", parsed.LookUpField("Password").DBName)
	assert.Equal(t, "banned", parsed.LookUpField("Banned").DBName)
	assert.Equal(t, "banned_reason", parsed.LookUpField("BannedReason").DBName)
	assert.Equal(t, "pub_path", parsed.LookUpField("PubPath").DBName)
	assert.Equal(t, "pub_pict", parsed.LookUpField("PubPict").DBName)
	assert.Equal(t, "date_registered", parsed.LookUpField("DateRegistered").DBName)
	assert.Equal(t, "last_update", parsed.LookUpField("LastUpdate").DBName)

	// email and username are two separate columns in account.users
	assert.NotEqual(t, parsed.LookUpField("Email").DBName, parsed.LookUpField("Username").DBName)
}
