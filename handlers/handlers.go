package handlers

import (
	"Pdf-Management/database"
	"Pdf-Management/models"
	"Pdf-Management/render"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
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
func Pdf(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "pdf.page.tmpl", &models.TemplateData{})
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
	http.Redirect(w, r, "/upload-pdf", http.StatusSeeOther)
}

func StorePDF(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	email := session.Values["email"].(string)
	w.Write([]byte("Email from the session is " + email))

	err := r.ParseMultipartForm(32 << 20) // Set max size for uploaded files
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the uploaded file from the form
	file, _, err := r.FormFile("pdfFile")
	title := r.FormValue("title")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Insert the PDF file into the database
	err = database.InsertPdf(email, file, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect or respond with a success message
	//http.Redirect(w, r, "/success", http.StatusFound)
	w.Write([]byte("Pdf File Submitted successfully"))
}

func GetPDF(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	email := session.Values["email"].(string)
	if email == "" {
		http.Error(w, "Email parameter is missing", http.StatusBadRequest)
		return
	}

	pdfData, err := database.GetPdf(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the appropriate headers for PDF file response
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=example.pdf")

	// Write the PDF data to the response writer
	_, err = w.Write(pdfData)
	if err != nil {
		log.Println("Failed to write PDF data to response:", err)
	}

}



func GetAllPdf(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	email := session.Values["email"].(string)
	if email == "" {
		http.Error(w, "Email parameter is missing", http.StatusBadRequest)
		return
	}
	pdfDataList, err := database.GetAllUserPdf(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	
	htmlTemplate := `

	<!DOCTYPE html>
<html>
<head>
	<title>PDF Viewer</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			background-color: #f8f8f8;
			padding: 20px;
		}

		.container {
			max-width: 600px;
			margin: 0 auto;
			background-color: #fff;
			border-radius: 8px;
			box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
			padding: 20px;
		}

		.pdf-list {
			list-style: none;
			padding: 0;
			margin: 0;
		}

		.pdf-item {
			margin-bottom: 10px;
		}

		.pdf-button {
			display: inline-block;
			padding: 10px 15px;
			border-radius: 4px;
			border: none;
			background-color: #007bff;
			color: #fff;
			font-size: 16px;
			cursor: pointer;
		}

		.pdf-button:hover {
			background-color: #0056b3;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>PDF Viewer</h1>
		<ul class="pdf-list">
			{{range $index, $pdfData := .}}
				<li class="pdf-item">
					<form action="/pdf" method="post" target="_blank">
						<input type="hidden" name="pdfData" value="{{base64Encode $pdfData.Data}}">
						<button class="pdf-button" type="submit">View PDF {{$index}} - {{$pdfData.Title}}</button>
					</form>
				</li>
			{{end}}
		</ul>
	</div>
</body>
</html>
`

	// Create a template function for base64 encoding
	funcMap := template.FuncMap{
		"base64Encode": func(data []byte) string {
			return base64.StdEncoding.EncodeToString(data)
		},
	}

	// Parse the HTML template
	tmpl, err := template.New("pdfs").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template and pass the PDF data list
	err = tmpl.Execute(w, pdfDataList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandlePDF(w http.ResponseWriter, r *http.Request) {
	pdfData := r.FormValue("pdfData")
	if pdfData == "" {
		http.Error(w, "PDF data is missing", http.StatusBadRequest)
		return
	}

	decodedData, err := base64.StdEncoding.DecodeString(pdfData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the appropriate headers for PDF file response
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline")

	// Write the PDF data to the response writer
	_, err = w.Write(decodedData)
	if err != nil {
		log.Println("Failed to write PDF data to response:", err)
	}

}
