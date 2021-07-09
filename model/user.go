package model

import "time"

type User struct {
	Id       int    `json:"id"  `
	Email    string `json:"email"`   // length 50
	Password string `json:"password"` // length 255
	Name     string `json:"name"`    // length 50

	Phone    string `json:"phone"`
	Address  string `json:"address"`

	Role      int       `json:"role"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
