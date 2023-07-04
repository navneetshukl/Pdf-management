package handlers

import (
	"Pdf-Management/database"
	"Pdf-Management/models"
	"Pdf-Management/render"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "signup.page.tmpl", &models.TemplateData{})
}
func Login(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "login.page.tmpl", &models.TemplateData{})
}
func Pdf(w http.ResponseWriter,r* http.Request){
	render.RenderTemplate(w,"pdf.page.tmpl",&models.TemplateData{})
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var name, email, password string
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
	if userEmail == "" {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	fmt.Println("UserPassword is ", userPassword)

	if err != nil {
		// Passwords don't match
		fmt.Println("Password Error ", err)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "session-name")
	session.Values["email"] = email
	session.Save(r, w)

	//w.Write([]byte(userEmail))

	fmt.Println("User Authenticated")
	http.Redirect(w,r,"/pdf",http.StatusSeeOther)
}

func StorePDF(w http.ResponseWriter,r* http.Request){
	session, _ := store.Get(r, "session-name")
	email := session.Values["email"].(string)
	w.Write([]byte("Email from the session is " +email))

	err := r.ParseMultipartForm(32 << 20) // Set max size for uploaded files
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retrieve the uploaded file from the form
		file, _, err := r.FormFile("pdfFile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Insert the PDF file into the database
		err = database.InsertPdf(email, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect or respond with a success message
		//http.Redirect(w, r, "/success", http.StatusFound)
		w.Write([]byte("Pdf File Submitted successfully"))
}
