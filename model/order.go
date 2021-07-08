package model

import "time"

type Order struct {
	Id          int     `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Phone       string  `json:"phone"`
	Address     string  `json:"add"`
	Email       string  `json:"email"`
	Note        string  `json:"note"`
	Total       float64 `json:"total"`
	TotalBill   float64 `json:"total_bill"`
	TotalWeight string  `json:"total_weight"`
	Shipping    float64 `json:"shipping"`

	VoucherCode        string `json:"voucher_code"`
	PaymentMethod      string `json:"payment_method"`
	SellerNote         string `json:"seller_note"`
	MailDeliveryStatus int    `json:"mail_delivery_status"`
	UserId             int    `json:"user_id"`

	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
