package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/o7z/try_gorm/bussiness"

	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		"root", "123", "127.0.0.1", 3306, "try_gorm")
	if db, err := gorm.Open("mysql", connStr); err != nil {
		fmt.Printf("gorm.Open(`mysql`,`%s`)->%s\n\r", connStr, err)
	} else {
		defer db.Close()
		if err := bussiness.Init(db); err != nil {
			fmt.Printf("bussiness.Init(db)->%s\n\r", err)
		} else if err := bussiness.Test(); err != nil {
			fmt.Printf("bussiness.Test()->%s\n\r", err)
		}
	}
}
