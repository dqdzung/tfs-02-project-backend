package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"project-backend/database"
	"project-backend/model"
	"project-backend/order/request"
	"project-backend/util/message"
	response "project-backend/util/response"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var db *gorm.DB

func CheckCart(w http.ResponseWriter, r *http.Request)  {
	db = database.ConnectDB()
//	1 Get user ? Sau sửa sau khi login sẽ gửi thông tin người dùng và lưu vào state
//	emailUser := r.Header.Get("email")
//	user := model.User{}
//	err := user.GetByEmail(emailUser)
//	if err != nil {
//		response.RespondWithJSON(w, 401, 0, "User not exists", nil)
//		return
//	}
// 2. get request
	requestCart := request.RequestCheckCart{}
	err = json.NewDecoder(r.Body).Decode(&requestCart)
	if err != nil {
		response.RespondWithJSON(w, 400, 0, message.ERROR_BAD_REQUEST, nil)
		return
	}
// 3. Check voucher
	discount := 0.0
	unit := "usd"
	maxSaleAmount := 0.0
	if requestCart.VoucherCode != "" {
		voucher := model.Voucher{}
		err = voucher.GetByCode(requestCart.VoucherCode)
		if err != nil {
			response.RespondWithJSON(w, 400, 0, message.ERROR_VOUCHER_NOT_EXISTS, nil)
			return
		}
		if !time.Now().Before(voucher.TimeEnd) {
			response.RespondWithJSON(w, 400, 0, message.ERROR_VOUCHER_EXPIRED, nil)
			return
		}
		discount = voucher.Discount
		unit = voucher.Unit
		maxSaleAmount = voucher.MaxSaleAmount
	}
// check variant
	for _, item := range requestCart.Cart {
		err = checkItem(db,&item)
		if err != nil {
			response.RespondWithJSON(w,400,0,err.Error(),item)
			return
		}
	}
	discountAmount := request.CaculateDiscountAmount(requestCart.Total,discount,maxSaleAmount,unit)
// check total
	err = requestCart.CheckCaculation(discountAmount)
	if err != nil {
		response.RespondWithJSON(w,400,0,err.Error(),nil)
		return
	}
	response.RespondWithJSON(w,200,1,message.SUCCESS,nil)
}


func checkItem(db *gorm.DB,item *request.ItemCheckCart) error{
	//check price, quantity, variant exist?
	sql := "SELECT quantity FROM variants WHERE id = ?  AND price = ? AND product_id = ?"
	quantity := 0
	db.Raw(sql, item.Variant.Id, item.Variant.Price, item.Id).Scan(&quantity)
	if item.Quantity > int64(quantity) {
		return errors.New(message.ERROR_PRODUCT_CHANGED)
	}
	return nil
}
func GetVoucherByCode(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	code := param["code"]
	voucher := model.Voucher{}
	err := voucher.GetByCode(code)
	
	if err != nil {
		response.RespondWithJSON(w, 400, 0, message.ERROR_VOUCHER_NOT_EXISTS, nil)
		return
	}
	if !time.Now().Before(voucher.TimeEnd) {
		response.RespondWithJSON(w, 400, 0, message.ERROR_VOUCHER_EXPIRED, nil)
		return
	}
	response.RespondWithJSON(w, 200, 1, "", voucher)

}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	//1. Get user
	db := database.ConnectDB()
	emailUser := r.Header.Get("email")
	user := model.User{}
	err := user.GetByEmail(emailUser)
	if err != nil {
		response.RespondWithJSON(w, 400, 0, "User not exists", nil)
		return
	}
	// get request order
	requestOrder := request.RequestCreateOrder{}
	err = json.NewDecoder(r.Body).Decode(&requestOrder)
	if err != nil {
		response.RespondWithJSON(w, 400, 0, "Bad request", nil)
		return
	}

	// response.RespondWithJSON(w, 200, 1, "", user)
	//2. Valid voucher code
	discount := 0.0
	unit := "usd"
	maxSaleAmount := 0.0
	voucher := model.Voucher{}

	err = voucher.GetByCode(requestOrder.VoucherCode)
	if err != nil {
		response.RespondWithJSON(w, 400, 0, "Voucher not exists", nil)
		return
	}
	if !time.Now().Before(voucher.TimeEnd) {
		response.RespondWithJSON(w, 400, 0, "Voucher expired", nil)
		return
	}
	discount = voucher.Discount
	unit = voucher.Unit
	maxSaleAmount = voucher.MaxSaleAmount

	// valid request
	err = requestOrder.ValidRequestCreateOrder()
	if err != nil {
		response.RespondWithJSON(w, 400, 0, err.Error(), nil)
		return
	}
	// check discountAmount
	err = requestOrder.CheckDiscountAmount(discount, unit, maxSaleAmount)
	if err != nil {
		response.RespondWithJSON(w, 400, 0, err.Error(), nil)
		return
	}
	// check total
	err = requestOrder.CheckTotal()
	if err != nil {
		response.RespondWithJSON(w, 400, 0, err.Error(), nil)
		return
	}
	// check price, weight request vs db
	sql := "SELECT * FROM variants WHERE id = ?  AND price = ? AND weight = ? AND product_id = ?"
	var resutlQuery *gorm.DB
	var variant model.Variant
	for _, item := range requestOrder.Carts {
		resutlQuery = db.Raw(sql, item.Id, item.Price, item.Weight, item.ProductId).Take(&variant)
		if resutlQuery.RowsAffected < 1 {
			response.RespondWithJSON(w, 400, 0, "Item "+strconv.Itoa(int(item.Id))+" not exists", nil)
			return
		}
	}
	// luu db
	tx := db.Begin()
	var orderDB model.Order
	orderDB.MapFromCreateOrder(&requestOrder, int64(user.Id))
	err = tx.Create(&orderDB).Error
	if err != nil {
		tx.Rollback()
		response.RespondWithJSON(w, 500, 0, "create err: "+err.Error(), nil)
		return
	}
	// update quantity
	for _, i := range requestOrder.Carts {
		sql = "UPDATE variants v SET v.quantity = v.quantity- ? WHERE v.id = ?"
		err = tx.Exec(sql, i.Quantity, i.Id).Error
		if err != nil {
			tx.Rollback()
			response.RespondWithJSON(w, 500, 0, "The variant is not in stock", i)
			return
		}
		sql = "UPDATE products p SET p.quantity = p.quantity- ? WHERE p.id = (SELECT v.product_id FROM variants  v WHERE v.id = ? )"
		err = tx.Exec(sql, i.Quantity, i.Id).Error
		if err != nil {
			tx.Rollback()
			response.RespondWithJSON(w, 500, 0, "The product is not in stock", i)
			return
		}
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		response.RespondWithJSON(w, 500, 0, "create err: "+err.Error(), nil)
		return
	}
	response.RespondWithJSON(w, 201, 1, "success", nil)

}
