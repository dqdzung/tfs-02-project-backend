package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"project-backend/database"
	"project-backend/model"
	"project-backend/order/request"
	"project-backend/util/constant"
	response "project-backend/util/response"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// var db *gorm.DB

func CheckCart(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()
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
	err := json.NewDecoder(r.Body).Decode(&requestCart)
	if err != nil {
		response.RespondWithJSON(w, 400, 0, constant.ERROR_BAD_REQUEST, nil)
		return
	}
	if len(requestCart.Cart) == 0{
		response.RespondWithJSON(w, 400, 0, constant.ERROR_BAD_REQUEST, nil)
		return
	}
	// 3. Check voucher
	discount := 0.0
	unit := constant.UNIT_USA
	maxSaleAmount := 0.0
	if requestCart.VoucherCode != "" {
		voucher := model.Voucher{}
		err = voucher.GetByCode(requestCart.VoucherCode)
		if err != nil {
			response.RespondWithJSON(w, 400, 0, constant.ERROR_VOUCHER_NOT_EXISTS, nil)
			return
		}
		if !time.Now().Before(voucher.TimeEnd) {
			response.RespondWithJSON(w, 400, 0, constant.ERROR_VOUCHER_EXPIRED, nil)
			return
		}
		discount = voucher.Discount
		unit = voucher.Unit
		maxSaleAmount = voucher.MaxSaleAmount
	}
	// check variant
	for _, item := range requestCart.Cart {
		err = checkItem(db, &item)
		if err != nil {
			response.RespondWithJSON(w, 400, 0, err.Error(), item)
			return
		}
	}
	discountAmount := request.CaculateDiscountAmount(requestCart.Total, discount, maxSaleAmount, unit)
	// check total
	err = requestCart.CheckCaculation(discountAmount)
	if err != nil {
		response.RespondWithJSON(w, 400, 0, err.Error(), nil)
		return
	}
	response.RespondWithJSON(w, 200, 1, constant.SUCCESS, nil)
}

func checkItem(db *gorm.DB, item *request.ItemCheckCart) error {
	//check price, quantity, variant exist?
	sql := "SELECT quantity FROM variants WHERE id = ?  AND price = ? AND product_id = ?"
	quantity := 0
	db.Raw(sql, item.Variant.Id, item.Variant.Price, item.Id).Scan(&quantity)
	if item.Quantity > int64(quantity) {
		return errors.New(constant.ERROR_PRODUCT_CHANGED)
	}
	return nil
}
func GetVoucherByCode(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	code := param["code"]
	voucher := model.Voucher{}
	err := voucher.GetByCode(code)

	if err != nil {
		response.RespondWithJSON(w, 400, 0, constant.ERROR_VOUCHER_NOT_EXISTS, nil)
		return
	}
	if !time.Now().Before(voucher.TimeEnd) {
		response.RespondWithJSON(w, 400, 0, constant.ERROR_VOUCHER_EXPIRED, nil)
		return
	}
	response.RespondWithJSON(w, 200, 1, "", voucher)

}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	// //1. Get user
	db := database.ConnectDB()
	emailUser := r.Header.Get("email")
	user := model.User{}
	err := user.GetByEmail(emailUser)
	if err != nil {
		response.RespondWithJSON(w, 401, 0, "User not exists", nil)
		return
	}
	// get request order
	requestOrder := request.RequestCreateOrder{}
	err = json.NewDecoder(r.Body).Decode(&requestOrder)
	if err != nil {
		response.RespondWithJSON(w, 400, 0, constant.ERROR_BAD_REQUEST, nil)
		return
	}
	if len(requestOrder.Cart) == 0{
		response.RespondWithJSON(w, 400, 0, constant.ERROR_BAD_REQUEST, nil)
		return
	}
	// valid customer information
	err = requestOrder.ValidCustomerInformation()
	if err != nil {
		response.RespondWithJSON(w, 400, 0, err.Error(), nil)
		return
	}

	//Valid voucher code
	discount := 0.0
	unit := constant.UNIT_USA
	maxSaleAmount := 0.0
	if requestOrder.VoucherCode != "" {
		voucher := model.Voucher{}
		err = voucher.CheckVoucher(requestOrder.VoucherCode)
		if err != nil {
			response.RespondWithJSON(w, 400, 0, err.Error(), nil)
			return
		}
		discount = voucher.Discount
		unit = voucher.Unit
		maxSaleAmount = voucher.MaxSaleAmount
	}


	//check item
	for _, item := range requestOrder.Cart {
		err = checkItem(db, &item)
		if err != nil {
			response.RespondWithJSON(w, 400, 0, err.Error(), item)
			return
		}
	}
	// check discountAmount
	discountAmount := request.CaculateDiscountAmount(requestOrder.Total, discount, maxSaleAmount, unit)
	// check total
	err = requestOrder.CheckCaculation(discountAmount)
	if err != nil {
		response.RespondWithJSON(w, 400, 0, err.Error(), nil)
		return
	}

	// luu db
	tx := db.Begin()
	var orderDB model.Order
	orderDB.MapFromCreateOrder(&requestOrder, int64(user.Id))
	err = tx.Create(&orderDB).Error
	if err != nil {
		tx.Rollback()
		response.RespondWithJSON(w, 500, 0, constant.ERROR_SERVER, nil)
		return
	}
	// update quantity
	for _, i := range requestOrder.Cart {
		sql := "UPDATE variants v SET v.quantity = v.quantity- ? WHERE v.id = ?"
		err = tx.Exec(sql, i.Quantity, i.Variant.Id).Error
		if err != nil {
			tx.Rollback()
			response.RespondWithJSON(w, 500, 0, constant.ERROR_SERVER, i)
			return
		}
		sql = "UPDATE products p SET p.quantity = p.quantity- ? WHERE p.id = (SELECT v.product_id FROM variants  v WHERE v.id = ? )"
		err = tx.Exec(sql, i.Quantity, i.Variant.Id).Error
		if err != nil {
			tx.Rollback()
			response.RespondWithJSON(w, 500, 0, constant.ERROR_SERVER, i)
			return
		}
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		response.RespondWithJSON(w, 500, 0, constant.ERROR_SERVER, nil)
		return
	}
	response.RespondWithJSON(w, 201, 1, constant.SUCCESS, nil)

}
