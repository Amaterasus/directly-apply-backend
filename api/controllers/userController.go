package controllers

import (
	"net/http"

	"encoding/json"

	"github.com/Amaterasus/direct-apply-backend/api/models"
)

type formData struct {
	Name string `json:"name"`
	Email string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
	FoundUs string `json:"foundUs"`
	SendJobMatches bool `json:"sendJobMatches"`
	AgreedToTerms bool `json:"agreedToTerms"`
}

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

func login(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	password := r.FormValue("password")

	user := models.User{}

	if user.Authorise(name, password) {
		json.NewEncoder(w).Encode(user)
	} else {
		m := make(map[string]string)
		m["Message"] = "name and password do not match"
		json.NewEncoder(w).Encode(m)
	}
}

func wakeup(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]string)
	m["Message"] = "Server is now awake"
	json.NewEncoder(w).Encode(m)
}

func signUp(w http.ResponseWriter, r *http.Request) {

	var data formData

	json.NewDecoder(r.Body).Decode(&data)

	user := &models.User{}

	if data.Password == data.PasswordConfirmation && data.AgreedToTerms {
		newUser := user.Create(data.Name, data.Email, data.PhoneNumber, data.Password, data.FoundUs, data.SendJobMatches, data.AgreedToTerms)

		json.NewEncoder(w).Encode(newUser)
	} else {
		m := make(map[string]string)
		m["Error"] = "Password does not match!"
		json.NewEncoder(w).Encode(m)
	}

}
