package handlers

import (
	"Pdf-Management/database"
	"Pdf-Management/models"
	"Pdf-Management/render"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "signup.page.tmpl", &models.TemplateData{})
}
func Login(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "login.page.tmpl", &models.TemplateData{})
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var name,email,password string
	name = r.FormValue("name")
	email = r.FormValue("email")
	password = r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error in generating the Hashed Password")
	}
	database.InsertUser(name, email, string(hashedPassword))
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	var email, password string
	email = r.FormValue("email")
	password = r.FormValue("password")
	fmt.Println(email)
	fmt.Println(password)

	var userEmail, userPassword string
	var err error

	userEmail, userPassword, err = database.GetUser(email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	if userEmail==""{
		http.Error(w, "User not found", http.StatusUnauthorized)
		return;
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	fmt.Println("UserPassword is ", userPassword)

	if err != nil {
		// Passwords don't match
		fmt.Println("Password Error ", err)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	w.Write([]byte(userEmail))

	fmt.Println("User Authenticated")
	// Continue with further actions or return a success response
}