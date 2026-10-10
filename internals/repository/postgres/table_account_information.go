package internal_repository_postgres

import "time"

const (
	TABLE_ACCOUNT_INFORMATION_FILE = "migrations/000003_table_account_information.sql"
)

// Maps to account.informations
type UsersInformationRow struct {
	IDUser       string     `gorm:"column:id_user"`
	FirstName    *string    `gorm:"column:first_name"`
	MiddleName   *string    `gorm:"column:midle_name"`
	LastName     *string    `gorm:"column:last_name"`
	Street1      *string    `gorm:"column:street1"`
	Street2      *string    `gorm:"column:street2"`
	District     *string    `gorm:"column:district"`
	City         *string    `gorm:"column:city"`
	Province     *string    `gorm:"column:province"`
	PostalCode   *string    `gorm:"column:postal_code"`
	PhoneNumbers []*string  `gorm:"column:phone_numbers"`
	BackupEmails []*string  `gorm:"column:backup_emails"`
	LastUpdate   *time.Time `gorm:"column:last_update"`
}

func (UsersInformationRow) TableName() string {
	return "account.informations"
}
