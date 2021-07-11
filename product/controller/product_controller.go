package controller

import (
	"math"
	"net/http"
	"project-backend/database"

	response "project-backend/product/response"
	resultResponse "project-backend/util/response"

	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var db = database.ConnectDB()

const (
	DEFAULT_PAGE_SIZE  = 10
	DEFAULT_PAGE_INDEX = 1
	DEFAULT_SORT_TITLE = "id"
	DEFAULT_SORT_BY    = "asc"
)

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	productId, err := strconv.Atoi(param["id"])
	if err != nil {
		resultResponse.RespondWithJSON(w, 400, 0, "Invalid param", "")
		return
	}
	var result response.ResponseProductByID
	sql := "SELECT * FROM products WHERE active = 2 AND  id = ?"
	resultQuery := db.Raw(sql, productId).Scan(&result)
	if resultQuery.RowsAffected < 1 {
		resultResponse.RespondWithJSON(w, 404, 1, "Not found", "")
		return
	}

	// get các options, option_value
	var options []response.Option
	db.Preload("OptionValues").Find(&options)
	result.SetOptions(&options)

	// get variants
	var variants []response.Variant
	sql = "SELECT * FROM `variants` WHERE product_id = ?"
	db.Raw(sql, productId).Scan(&variants)
	result.SetVariants(&variants)

	//get image
	var images []response.Image
	sql = "SELECT * FROM product_images WHERE product_id = ?"
	db.Raw(sql, productId).Scan(&images)
	result.SetImages(&images)

	resultResponse.RespondWithJSON(w, 200, 1, "", result)
}
func SearchProduct(w http.ResponseWriter, r *http.Request) {
	//  input/: param:
	// pageindex: default: 1
	// name: product(like) // defaullt: ""
	// fillter:
	//brand
	//category // default: all
	// sort title : price | created time  // default: id
	// sort by: asc | desc	// default: asc
	//---------------------------------------------

	// lấy param
	nameProduct := r.URL.Query().Get("name") // nameProduct theo tên product
	brand := r.URL.Query().Get("brand")
	category := r.URL.Query().Get("category")
	sortTitle := r.URL.Query().Get("sort")
	orderBy := r.URL.Query().Get("order")
	rawPageIndex := r.URL.Query().Get("page")
	rawLimit := r.URL.Query().Get("limit")

	// check param: số, ký tự đặc biệt
	nameProduct = strings.Replace(nameProduct, "%", "", -1)
	nameProduct = strings.Replace(nameProduct, "-", "", -1)

	pageIndex, err := strconv.Atoi(rawPageIndex)
	if err != nil {
		pageIndex = DEFAULT_PAGE_INDEX
	}

	// query sql
	sql := "active = 2 AND products.name LIKE '%" + nameProduct + "%' "
	if brand != "" {
		sql += "AND products.brand_name = '" + brand + "' "
	}
	if category != "" {
		sql += "AND products.category_id = (SELECT id FROM categories WHERE categories.name = '" + category + "') "
	}
	switch sortTitle {
	case "price":
	case "date":
		sortTitle = "created_at"
	default:
		sortTitle = DEFAULT_SORT_TITLE
	}
	if orderBy != "asc" && orderBy != "desc" {
		orderBy = DEFAULT_SORT_BY
	}
	sql += "ORDER BY " + sortTitle + " " + orderBy + " "
	limit, err := strconv.Atoi(rawLimit)
	if err != nil {
		limit = DEFAULT_PAGE_SIZE
	}
	if limit < 1 {
		limit = DEFAULT_PAGE_SIZE
	}
	if pageIndex < 1 {
		pageIndex = DEFAULT_PAGE_INDEX
	}
	// check page total
	var totalElement int64
	db.Raw("SELECT COUNT(*) FROM products WHERE " + sql).Scan(&totalElement)
	if totalElement == 0 {
		resultResponse.RespondWithJSON(w, 200, 1, "Not found", nil)
		return
	}
	totalPage := int(math.Ceil(float64(totalElement) / float64(limit)))

	if pageIndex < 1 {
		pageIndex = DEFAULT_PAGE_INDEX
	}
	if pageIndex > totalPage {
		pageIndex = totalPage
	}
	// sql thu được
	//	SELECT id, name, image_url price,original_price,quantity FROM products WHERE active = 2
	//	products.name LIKE '% %' AND
	//	products.brand_name = '' AND
	//	products.category_id = (SELECT id FROM categories WHERE categories.name = 'dog')
	//	ORDER BY produc ASC
	//	LIMIT 0,10
	sql = "SELECT * FROM products WHERE " + sql + " LIMIT " + strconv.Itoa((pageIndex-1)*limit) + " , " + strconv.Itoa(limit)
	// Lấy data từ dât
	var products []response.Product
	db.Raw(sql).Scan(&products)

	res := response.ResponseSearch{
		TotalPage:    totalPage,
		TotalElement: int(totalElement),
		PageIndex:    pageIndex,
		PageSize:     limit,
		Products:     products,
	}
	resultResponse.RespondWithJSON(w, 200, 1, "", res)

}
