package internal_repository_postgres

import "time"

const (
	TABLE_ACCOUNT_USERS_FILE = "migrations/000002_table_account_users.sql"
)

// Maps account.users
type UsersRow struct {
	ID             string     `gorm:"column:id;primaryKey"`
	IDRef          string     `gorm:"column:id_ref"`
	Email          string     `gorm:"column:email"`
	PhoneNumber    string     `gorm:"column:phone_number"`
	Username       string     `gorm:"column:email"`
	Password       string     `gorm:"column:password"`
	Banned         bool       `gorm:"column:banned"`
	BannedReason   *string    `gorm:"column:banned_reason"`
	PubPath        string     `gorm:"column:pub_path"`
	PubPict        string     `gorm:"column:pub_pict"`
	DateRegistered time.Time  `gorm:"column:date_registered"`
	LastUpdate     *time.Time `gorm:"column:last_update"`
}

func (UsersRow) TableName() string {
	return "account.users"
}
