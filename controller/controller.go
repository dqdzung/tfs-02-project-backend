package controller

import (
	"encoding/json"
	"net/http"
	"project-backend/database"
	"strconv"

	"github.com/gorilla/mux"
)

var db = database.ConnectDB()

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products := []database.Product{}

	result := db.Find(&products)

	if result.Error != nil {
		res := map[string]interface{}{
			"success": 0,
			"message": result.Error,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := map[string]interface{}{
		"success": 1,
		"data":    &products,
	}
	json.NewEncoder(w).Encode(res)
}

func GetOneProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	product := database.Product{}
	result := db.First(&product, id)

	if result.Error != nil {
		res := map[string]interface{}{
			"success": 0,
			"message": result.Error,
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	res := map[string]interface{}{
		"success": 1,
		"data":    &product,
	}
	json.NewEncoder(w).Encode(res)
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	newProduct := database.Product{}

	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		res := map[string]interface{}{
			"success": 0,
			"message": err,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	if result := db.Create(&newProduct); result.Error != nil {
		res := map[string]interface{}{
			"success": 1,
			"message": result.Error,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := map[string]interface{}{
		"success": 1,
		"data":    &newProduct,
	}
	json.NewEncoder(w).Encode(res)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// vars := mux.Vars(r)
	// id, _ := strconv.Atoi(vars["id"])

	// student := database.Student{}
	// result := db.First(&student, id)
	// if result.Error != nil {
	// 	fmt.Fprintf(w, "No entry at id %v", id)
	// 	return
	// }

	// if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
	// 	fmt.Fprintf(w, "error when parsing body %v", err)
	// 	return
	// }

	// db.Save(&student)

	// fmt.Fprintf(w, "Updated id %v to %v", id, student)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// vars := mux.Vars(r)
	// id, _ := strconv.Atoi(vars["id"])

	// student := database.Student{}
	// result := db.Delete(&student, id)
	// if result.Error != nil {
	// 	fmt.Fprintf(w, "No entry at id %v", id)
	// 	return
	// }
	// fmt.Fprintf(w, "Deleted id %v", id)
}

// GET ALL
func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

// POST ONE
func AddOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

// GET ONE
func GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

// PUT ONE
func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

// DELETE ONE
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
