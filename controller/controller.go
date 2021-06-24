package controller

import (
	"net/http"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// var students = []database.Student{}
	// results := db.Find(&students)
	// json.NewEncoder(w).Encode(results.Value)
}

func GetOneProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// vars := mux.Vars(r)
	// id, _ := strconv.Atoi(vars["id"])

	// student := database.Student{}
	// result := db.First(&student, id)

	// if result.Error != nil {
	// 	fmt.Fprintf(w, "No entry at id %v", id)
	// 	return
	// }

	// json.NewEncoder(w).Encode(result.Value)
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// newStudent := database.Student{}

	// if err := json.NewDecoder(r.Body).Decode(&newStudent); err != nil {
	// 	fmt.Fprintf(w, "error when parsing body %v", err)
	// 	return
	// }

	// if result := db.Create(&newStudent); result.Error != nil {
	// 	fmt.Fprintf(w, "Couldn't add %v. Error: %v", newStudent, result.Error)
	// 	return
	// }

	// fmt.Fprintf(w, "Added %v", newStudent)
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
