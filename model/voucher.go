package model

import "time"

type Voucher struct {
	Id            int       `json:"id"`
	Code          string    `json:"code"`
	Discount      float64   `json:"discount"` // giam
	Unit          string    `json:"unit"`  // persent || usd
	MaxSaleAmount float64   `json:"max_sale_amount"`
	Description   string    `json:"description"`
	TimeEnd       time.Time `json:"time_end"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
func (v *Voucher) GetVoucherByCode(code string) error{
	return 	db.Where("code = ?", code).Take(v).Error
}
