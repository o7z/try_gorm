package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	ID         string `gorm:"size:50"`
	Name       string
	Articles   []Article `gorm:"foreignkey:owner"`
	CreateTime time.Time
}

type Article struct {
	ID         string `gorm:"size:50"`
	Title      string
	Content    string
	CreateUser string `gorm:"size:50"`
	CreateTime time.Time
}

func main() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		"root", "123", "127.0.0.1", 3306, "try_gorm")
	if db, err := gorm.Open("mysql", connStr); err != nil {
		fmt.Printf("gorm.Open(`mysql`,`%s`)->%s\n\r", connStr, err)
	} else {
		defer db.Close()
		db.LogMode(true) //关心一下日志
		db.SingularTable(true)
		// db.CreateTable(&User{})
		user := &User{}
		article := &Article{}
		db.AutoMigrate(user)
		db.AutoMigrate(card)
		db.Model(card).AddForeignKey("owner", "user(id)", "RESTRICT", "RESTRICT")
		// db.AutoMigrate(&User{})
		// fmt.Printf("database connected->%v\n\r", db)
		// if db.HasTable(&User{}) {
		// 	fmt.Println("user exists")
		// } else {
		// 	fmt.Println("user does not exist")
		// 	// db.AutoMigrate(&User{})
		// 	db.CreateTable(&User{})
		// 	fmt.Println("check again : ", db.HasTable("users"))
		// }
	}
}
