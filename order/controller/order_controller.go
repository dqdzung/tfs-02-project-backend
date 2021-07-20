package controller

import (
	"encoding/json"
	"net/http"
	"project-backend/database"
	"project-backend/model"
	"project-backend/order/request"
	response "project-backend/util/response"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetVoucherByCode(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	code := param["code"]
	voucher := model.Voucher{}
	err := voucher.GetByCode(code)
	
	if err != nil {
		response.RespondWithJSON(w, 400, 0, "Voucher not exists", nil)
		return
	}
	if !time.Now().Before(voucher.TimeEnd) {
		response.RespondWithJSON(w, 400, 0, "Voucher expired", nil)
		return
	}
	response.RespondWithJSON(w, 200, 1, "", voucher)

}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	//1. Get user
	db := database.ConnectDB()
	emailUser := r.Header.Get("email")
	user := model.User{}
	err := user.GetUserByEmail(emailUser)
	if err != nil {
		response.RespondWithJSON(w, 400, 0, "User not exists", nil)
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
	sql := "SELECT * FROM variants WHERE id = ?  AND price = ? AND weight = ? AND alias = ?"
	var resutlQuery *gorm.DB
	var variant model.Variant
	for _, item := range requestOrder.Carts {
		resutlQuery = db.Raw(sql, item.Id, item.Price, item.Weight, item.Alias).Take(&variant)
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
