package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (db *gorm.DB) {
	dsn := "root:admin@/project?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error when connect to db ", err)
		return
	}

	if err != nil {
		log.Fatal("error when auto migrate table ", err)
	}
	return db
}

func CreateTable() {
	db := ConnectDB()

	db.Debug().Migrator().DropTable(&User{}, &Category{}, &Brand{}, &Product{}, &Order{}, &ProductOrder{})

	db.AutoMigrate(&User{}, &Category{}, &Brand{}, &Product{}, &Order{}, &ProductOrder{})

	fmt.Println("Done...")
}
