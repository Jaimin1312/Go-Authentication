package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"package/database"
	"package/model"
	"package/utility"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	connection := database.GetDatabase()
	defer database.Closedatabase(connection)

	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		var err model.Error
		err = utility.SetError(err, "Error in reading body")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	var dbuser model.User
	connection.Where("email = ?", user.Email).First(&dbuser)

	//check email is alredy register or not
	if dbuser.Email != "" {
		var err model.Error
		err = utility.SetError(err, "Email already in use")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	user.Password, err = utility.GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}

	//insert user details in database
	connection.Create(&user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	connection := database.GetDatabase()
	defer database.Closedatabase(connection)

	var authdetails model.Authentication

	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		var err model.Error
		err = utility.SetError(err, "Error in reading body")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var authuser model.User
	connection.Where("email = 	?", authdetails.Email).First(&authuser)

	if authuser.Email == "" {
		var err model.Error
		err = utility.SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	check := utility.CheckPasswordHash(authdetails.Password, authuser.Password)

	if !check {
		var err model.Error
		err = utility.SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := utility.GenerateJWT(authuser.Email, authuser.Role)
	if err != nil {
		var err model.Error
		err = utility.SetError(err, "Failed to generate token")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var token model.Token
	token.Email = authuser.Email
	token.Role = authuser.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOME PUBLIC INDEX PAGE"))
}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("NOT AUTHORIZED"))
		return
	}
	w.Write([]byte("ADMIN INDEX PAGE"))
}

func UserIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "user" {
		w.Write([]byte("NOT AUTHORIZED"))
		return
	}
	w.Write([]byte("User INDEX PAGE"))
}
