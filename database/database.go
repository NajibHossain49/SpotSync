package database

import (
	"log"
	"time"

	"spotsync-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens the PostgreSQL connection, configures the connection pool,
// and runs auto-migrations for all models.
func Connect(databaseURL string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure the underlying connection pool. This is important in
	// production so we don't exhaust database connections under load.
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)                 // max simultaneous open connections
	sqlDB.SetMaxIdleConns(10)                 // connections kept idle in the pool
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // recycle connections periodically

	// Auto-migrate all tables.
	if err := db.AutoMigrate(
		&models.User{},
		&models.ParkingZone{},
		&models.Reservation{},
	); err != nil {
		log.Fatalf("Auto-migration failed: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return db
}
