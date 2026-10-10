package internal_repository_postgres

import (
	// "context"
	// "fmt"
	"time"

	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Carries the pool settings of one connection
type Options struct {
	DSN             string
	MaxOpenConns    int
	MadIdleConss    int
	ConnMaxLifeTime time.Duration
	LogLevel        logger.LogLevel
}

const (
	// postgres:18
	PostgreSQL_IMAGE_VERSION = "postgres:18"
)

// --------------------------------------------------------- //

// // GORM handle bound to one DSN
// func Open(ctx context.Context, options Options) (*gorm.DB, error) {
// 	db, err := gorm.Open(postgres.Open(options.DSN), &gorm.Config{
// 		Logger: logger.Default.LogMode(options.LogLevel),
// 		NowFunc: func() time.Time { return time.Now().UTC() },
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("postgres: open connection: %w", err)
// 	}
//
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return nil, fmt.Errorf("postgres: reach sql handle: %w", err)
// 	}
//
// 	sqlDB.SetMaxOpenConns(options.MaxOpenConns)
// 	sqlDB.SetMaxIdleConns(options.MadIdleConss)
// 	sqlDB.SetConnMaxLifetime(options.ConnMaxLifeTime)
//
// 	err = sqlDB.PingContext(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("postgres: ping: %w", err)
// 	}
//
// 	return db, nil
// }
//
// // Reports whether the connection still up
// func Ping(ctx context.Context, db *gorm.DB) error {
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return fmt.Errorf("postgres: reach sql handle: %w", err)
// 	}
//
// 	if err := sqlDB.PingContext(ctx); err != nil {
// 		return fmt.Errorf("postgres: ping: %w", err)
// 	}
//
// 	return nil
// }
//
// // Release the pool
// func Close(db *gorm.DB) error {
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return fmt.Errorf("postgres: reach sql handle: %w", err)
// 	}
//
// 	if err := sqlDB.Close(); err != nil {
// 		return fmt.Errorf("postgres: close pool: %w", err)
// 	}
//
// 	return nil
// }
