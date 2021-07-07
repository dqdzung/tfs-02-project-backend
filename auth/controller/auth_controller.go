package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-backend/database"
	"project-backend/database/model"

	bcrypt "project-backend/util/bcrypt"
	jwt "project-backend/util/jwt"
	response "project-backend/util/response"
)

const (
	ParsingError  = "Error parsing request body"
	UserExist     = "User already exists"
	NewUser       = "New user added"
	NoUser        = "User doesn't exist"
	WrongPassword = "Wrong password"
	LoginSuccess  = "Logged in"
	TokenErr      = "Error generating token"
)

var db = database.ConnectDB()

func SignUp(w http.ResponseWriter, r *http.Request) {
	newUser := model.User{}

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		response.RespondWithJSON(w, 400, 0, ParsingError, nil)
		return
	}

	result := db.First(&model.User{}, "email = ?", newUser.Email)

	if result.RowsAffected != 0 {
		response.RespondWithJSON(w, 400, 0, UserExist, nil)
		return
	}

	// Encrypt password before saving to db (will be moved to frontend later)
	hashPassword, err := bcrypt.HashPassword(newUser.Password)
	if err != nil {
		fmt.Println("error hasing password", err)
	}

	newUser.Password = hashPassword

	if result := db.Create(&newUser); result.Error != nil {
		response.RespondWithJSON(w, 400, 0, UserExist, nil)
		return
	}
	response.RespondWithJSON(w, 200, 1, NewUser, &newUser)
}

func Login(w http.ResponseWriter, r *http.Request) {
	credentials := model.User{} //hold user login credentials from request body
	user := model.User{}        // hold user data from db

	// Get login data from request
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		response.RespondWithJSON(w, 400, 0, ParsingError, nil)
		return
	}

	// Check if requested email exists in db
	result := db.Where("email = ?", credentials.Email).First(&user)
	if result.Error != nil {
		response.RespondWithJSON(w, 400, 0, NoUser, nil)
		return
	}

	match := bcrypt.CheckPasswordHash(credentials.Password, user.Password)

	if !match {
		response.RespondWithJSON(w, 401, 0, WrongPassword, nil)
		return
	}

	//Generate jwt
	tokenString, err := jwt.CreateToken(w, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.RespondWithJSON(w, 500, 0, TokenErr, nil)
	}

	response.RespondWithJSON(w, 200, 1, LoginSuccess, map[string]string{"token": tokenString})
}
