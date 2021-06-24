package main

import (
	"fmt"
	"project-backend/database"
)

func main() {
	db := database.ConnectDB()

	db.Migrator().DropTable(&database.User{}, &database.Category{}, &database.Brand{}, &database.Product{}, &database.Order{}, &database.ProductOrder{})

	db.AutoMigrate(&database.User{}, &database.Category{}, &database.Brand{}, &database.Product{}, &database.Order{}, &database.ProductOrder{})

	fmt.Println("Done...")
}
