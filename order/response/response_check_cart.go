package response

type ResponseCheckCart struct{
	Name           string  `json:"name"`
	Phone          string  `json:"phone"`
	Address        string  `json:"address"`
	Email          string  `json:"email"`
	Total          float64 `json:"total"`
	DiscountAmount float64 `json:"discount_amount"`
	Shipping       float64 `json:"shipping"`
	TotalBill      float64 `json:"total_bill"`
	TotalWeight    string  `json:"total_weight"`
	VoucherCode    string  `json:"voucher_code"`
	Cart           []ItemCheckCart  `json:"cart"`
}
type ItemCheckCart struct { // variant
	Id         int64   `json:"id"`
	Alias      string  `json:"alias"`
	Image      string  `json:"image"`
	ProdutName string  `json:"name"`
	Quantity   int64   `json:"quantity"`
	Stock      int64   `json:"stock"`
	Variant    Variant `json:"variant"`
}
type Variant struct {
	Id            int64   `json:"id"`
	Code          string  `json:"code"`
	VariantName   string  `json:"name"`
	Option1       int64   `json:"option1"`
	Option2       int64   `json:"option2"`
	Option3       int64   `json:"option3"`
	OriginalPrice float64 `json:"original_price"`
	Position      int64   `json:"position"`
	Price         float64 `json:"price"`
	Stock      int64   `json:"quantity"`
	Weight        string  `json:"weight"`
}