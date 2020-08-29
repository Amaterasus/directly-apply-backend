package controllers

import (
	"net/http"

	"encoding/json"

	"github.com/Amaterasus/direct-apply-backend/api/models"
)

func allUsers(w http.ResponseWriter, r *http.Request) {
	
	user := models.User{}

	users := user.GetAllUsers()

	json.NewEncoder(w).Encode(users)
}

func authorised(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Authorised"]
	if token != nil {
		id := models.DecodeJWT(token[0])


		user := models.User{}

		user.FindUserByID(id)
		json.NewEncoder(w).Encode(user)
	} else {
		m := make(map[string]string)
		m["Error"] = "Not Authorised"
		json.NewEncoder(w).Encode(m)
	}
}

func signUp(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	phoneNumber := r.FormValue("phoneNumber")
	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("passwordConfirmation")
	foundUs := r.FormValue("foundUs")
	sendJobMatches := r.FormValue("sendJobMatches") == "true"
	agreedToTerms := r.FormValue("agreedToTerms") == "true"

	user := &models.User{}

	if (password == passwordConfirmation) {
		newUser := user.Create(name, email, phoneNumber, password, foundUs, sendJobMatches, agreedToTerms)

		json.NewEncoder(w).Encode(newUser)
	} else {
		m := make(map[string]string)
		m["Error"] = "Password does not match!"
		json.NewEncoder(w).Encode(m)
	}
	
}