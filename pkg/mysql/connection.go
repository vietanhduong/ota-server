package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func InitializeDatabase(cfg Config) (*DB, error) {
	_db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: cfg.DSN(),
	}))
	return &DB{_db}, err
}
