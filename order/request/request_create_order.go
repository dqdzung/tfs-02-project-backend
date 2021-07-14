package request

import (
	"errors"
	"project-backend/model"
	"project-backend/util/validator"
	"strconv"
)

type RequestCreateOrder struct {
	// Code
	Name          string  `json:"name"`
	Phone         string  `json:"phone"`
	Address       string  `json:"address"`
	Email         string  `json:"email"`
	Note          string  `json:"note"`
	Total         float64 `json:"total"`
	DiscountAmount float64 `json:"discount_amount"`
	Shipping      float64 `json:"shipping"`
	TotalBill     float64 `json:"total_bill"`
	TotalWeight   string  `json:"total_weight"`
	VoucherCode   string  `json:"voucher_code"`
	PaymentMethod string  `json:"payment_method"`
	Carts         []Item  `json:"carts"`
}
type Item struct { // variant
	Id            int64   `json:"id"`
	ProductName   string  `json:"product_name"`
	VariantName   string  `json:"variant_name"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"original_price"`
	Quantity      int64   `json:"quantity"`
	Weight        string  `json:"weight"`
	ProductId     int64   `json:"product_id"`
}

func (c RequestCreateOrder) CheckTotal(v model.Voucher) error {
	//discount := v.Discount
	//uint := v.Unit
	//maxSaleAmount := v.MaxSaleAmount

	total := 0.0
	totalBill := c.Shipping
	for _, item := range c.Carts {
		total += float64(item.Quantity) * item.Price
	}
	//switch uint {
	//case "percent":
	//
	//case "usd":
	//default:
	//
	//}
	totalBill += total
	if total != c.Total {
		return errors.New("invalid total ")

	}
	if totalBill != c.TotalBill {
		return errors.New("invalid total bill ")
	}
	return nil
}

// check total weight?

func (c RequestCreateOrder) ValidRequestCreateOrder() error {
	// check name
	if !(validator.CheckLength(c.Name, 50) && validator.CheckName(c.Name)) {
		return errors.New("invalid name")
	}
	if !validator.CheckPhone(c.Phone) {
		return errors.New("invalid phone")
	}
	if !validator.CheckLength(c.Address, 255) {
		return errors.New("invalid address")
	}
	if !validator.CheckMail(c.Email) {
		return errors.New("invalid email ")
	}
	if !validator.CheckLength(c.Note, 255) {
		return errors.New("invalid note ")
	}
	if c.Total < 0 || c.Shipping < 0 || c.TotalBill < 0 {
		return errors.New("invalid calculate total ")
	}

	//Check payment method?

	//Check cart
	for _, item := range c.Carts {
		if err := item.validItem(); err != nil {
			return errors.New(strconv.Itoa(int(item.Id)) + ": " + err.Error())
		}
	}
	return nil
}
func (i Item) validItem() error {
	if !validator.CheckLength(i.ProductName, 255) {
		return errors.New("invalid productName ")
	}
	if !validator.CheckLength(i.VariantName, 255) {
		return errors.New("invalid variantName ")
	}

	if i.Price < 0 {
		return errors.New("invalid price ")
	}
	if i.OriginalPrice != 0 && i.Price > i.OriginalPrice {
		return errors.New("invalid originalPrice ")
	}
	if i.Quantity < 0 {
		return errors.New("invalid quantity ")
	}
	if !validator.CheckLength(i.Weight, 255) {
		return errors.New("invalid weight ")
	}
	return nil
}
