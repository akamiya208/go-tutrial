package main

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
    ID   string `gorm:"primary_key"`
    Name string
	CreatedAt time.Time
}

func main() {

	// 詳細は https://github.com/go-sql-driver/mysql#dsn-data-source-name を参照
	dsn := "local:password@tcp(mysql:3306)/go_tutrial?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

    // Migrate the schema
    db.AutoMigrate(&User{})

	// Create
    db.Create(&User{ID: "user1", Name: "user1"})
	db.Create(&User{ID: "user2", Name: "user2"})
	db.Create(&User{ID: "user3", Name: "user3"})
}