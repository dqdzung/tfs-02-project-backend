package database

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (db *gorm.DB) {
	dsn := "root:Polarbear1011@/mysql_db?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error when connect to db ", err)
		return
	}
	db.AutoMigrate(&User{}, &Category{}, &Brand{}, &Product{}, &Order{}, &ProductOrder{})
	if err != nil {
		log.Fatal("error when auto migrate table ", err)
	}
	return db
}

type Product struct {
	Id            int            `json:"id"`
	Name          string         `json:"name"`
	Price         float64        `json:"price"`
	Sale          int            `json:"sale"`
	Quantity      int            `json:"qty"`
	Weight        float64        `json:"weight"`
	Description   string         `json:"desc"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	IsRendered    bool           `json:"isRendered"`
	BrandId       int            `db:"brandId"`
	CategoryId    int            `db:"categoryId"`
	ProductOrders []ProductOrder `gorm:"foreignKey:ProductId"`
}

type Brand struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Products []Product `gorm:"foreignKey:BrandId"`
}

type Order struct {
	Id            int            `json:"id"`
	Name          string         `json:"name"`
	Phone         string         `json:"phone"`
	Address       string         `json:"add"`
	Email         string         `json:"email"`
	Note          string         `json:"note"`
	Total         float64        `json:"total"`
	IsDone        bool           `json:"isDone"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	IsRendered    bool           `json:"isRendered"`
	ProductOrders []ProductOrder `gorm:"foreignKey:OrderId"`
}

type ProductOrder struct {
	Id        int `json:"id"`
	ProductId int `db:"productId"`
	OrderId   int `db:"orderId"`
}

type Category struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Products []Product `gorm:"foreignKey:CategoryId"`
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"pw"`
	Email    string `json:"email"`
}
