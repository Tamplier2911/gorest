package service

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

// NewMySQL is used to create a new MySQL database client.
func (s *Service) NewMySQL() (*gorm.DB, error) {

	// configure mysql connection
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		s.Config.MySQLUser, s.Config.MySQLPass, s.Config.MySQLHost, s.Config.MySQLDatabase,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:      zapgorm2.New(s.Logger.Desugar().Named("MySQL")),
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to my sqlserver: %s", err)
	}

	// create database
	s.Logger.Info("initializing database")
	err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", s.Config.MySQLDatabase)).Error
	if err != nil {
		s.Logger.Fatalw("failed to create database", "err", err)
	}

	// automigrate global models
	// err = db.AutoMigrate(&models.User{})
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to run migrations: %s", err)
	// }

	return db, nil
}
