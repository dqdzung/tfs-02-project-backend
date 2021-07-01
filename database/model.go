package database

import "time"

type Product struct {
	Id            int            `json:"id" gorm:"int"`
	Name          string         `json:"name" gorm:"type:varchar(255)"`
	Price         float64        `json:"price"  gorm:"float"`
	Sale          int            `json:"sale" gorm:"type:int"`
	Quantity      int            `json:"qty" gorm:"type:int"`
	Weight        float64        `json:"weight"  gorm:"float"`
	Description   string         `json:"desc"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	IsRendered    bool           `json:"isRendered" gorm:"type:boolean"`
	BrandId       int            `json:"brandId" gorm:"type:int"`
	CategoryId    int            `json:"categoryId" gorm:"type:int"`
	ProductOrders []ProductOrder `gorm:"foreignKey:ProductId"`
}
type Brand struct {
	Id       int       `json:"id"`
	Name     string    `json:"name" gorm:"type:string; size:50"`
	Products []Product `gorm:"foreignKey:BrandId"`
}
type Order struct {
	Id            int            `json:"id"`
	Name          string         `json:"name" gorm:"type:varchar(50)"`
	Phone         string         `json:"phone" gorm:"type:varchar(11)"`
	Address       string         `json:"add" gorm:"type:varchar(100)"`
	Email         string         `json:"email" gorm:"type:varchar(50)"`
	Note          string         `json:"note" gorm:"type:varchar(100)"`
	Total         float64        `json:"total"  gorm:"type:float" `
	IsDone        bool           `json:"isDone" gorm:"type:boolean"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	IsRendered    bool           `json:"isRendered" gorm:"type:boolean"`
	ProductOrders []ProductOrder `gorm:"foreignKey:OrderId"`
}
type ProductOrder struct {
	Id        int `json:"id"`
	ProductId int `gorm:"type:int"`
	OrderId   int `gorm:"type:int"`
}
type Category struct {
	Id       int       `json:"id" `
	Name     string    `json:"name" gorm:"type:varchar(50)"`
	Products []Product `gorm:"foreignKey:CategoryId"`
}
type User struct {
	Id       int    `json:"id"  `
	Username string `json:"username"  gorm:"type:varchar(50)"`
	Password string `json:"pw"  gorm:"type:varchar(50)"`
	Email    string `json:"email" gorm:"type:varchar(50)"`
}
