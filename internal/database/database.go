package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	const retries = 5

	for i := range retries {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Database connection established")
			return db, nil
		}
		fmt.Printf("Could not connect to database (attempt %d/%d): %v\n", i+1, retries, err)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts", retries)
}
