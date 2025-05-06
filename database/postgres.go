package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBInit(DBHost, DBUser, DBPassword, DBName, DBPort, DBSslMode string,
) (*gorm.DB, error) {
	dsn := "host=" + DBHost + " user=" + DBUser + " password=" + DBPassword + " dbname=" + DBName + " port=" + DBPort + " sslmode=" + DBSslMode
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return DB, nil
}
