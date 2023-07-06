package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/joho/godotenv"
)

type PDF struct {
	Title string
	Data  []byte
	Share string
}

func DB_Connect() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connString := os.Getenv("DB_CONNECTION_STRING")
	if connString == "" {
		log.Fatal("DB_CONNECTION_STRING not found in .env file")
	}
	conn, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Unable to connect: %v\n", err))
		return nil, err
	}

	//defer conn.Close()

	log.Println("Connected to database")

	err = conn.Ping()
	if err != nil {
		log.Fatal("Cannot Ping the database")
		return nil, err
	}
	log.Println("pinged database")

	return conn, nil

}

func InsertUser(name, email, password string) {
	conn, err := DB_Connect()
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	query := `insert into users(name,email,password) values($1,$2,$3)`

	_, err = conn.Exec(query, name, email, password)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted a User")
}

func GetUser(email string) (string, string, error) {
	conn, err := DB_Connect()
	defer conn.Close()
	if err != nil {
		return "", "", err
	}

	var userEmail, userPassword string

	query := `SELECT email,password FROM users WHERE email = $1`
	rows, err := conn.Query(query, email)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
		return "", "", err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&userEmail, &userPassword)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
			return "", "", err
		}
	}
	return userEmail, userPassword, nil
}

func CheckUser(email string) bool {
	conn, err := DB_Connect()
	defer conn.Close()
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		return false
	}
	

	// Prepare the SQL statement
	query := `SELECT email FROM users WHERE email = $1`
	rows, err := conn.Query(query, email)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
		return false
	}
	defer rows.Close()

	// Execute the query and retrieve the result
	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
			return false
		}
	}
	return count>0
}


func InsertPdf(email string, pdfFile multipart.File, title string,share string) error {
	conn, err := DB_Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	// Read the PDF file data
	pdfData, err := ioutil.ReadAll(pdfFile)
	if err != nil {
		return err
	}

	query := `INSERT INTO files (email, file, title, share) VALUES ($1, $2, $3, $4)`

	_, err = conn.Exec(query, email, pdfData, title, share)
	if err != nil {
		return err
	}

	log.Println("Inserted a PDF")
	return nil
}

func GetPdf(share string) ([]byte, error) {
	conn, err := DB_Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := `SELECT file FROM files WHERE share = $1`
	row := conn.QueryRow(query, share)

	var pdfData []byte
	err = row.Scan(&pdfData)
	if err != nil {
		return nil, err
	}

	return pdfData, nil
}


func GetAllUserPdf(email string) ([]PDF, error) {
	conn, err := DB_Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := `SELECT title, file, share FROM files WHERE email = $1`
	rows, err := conn.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pdfList []PDF

	for rows.Next() {
		var pdf PDF
		err = rows.Scan(&pdf.Title, &pdf.Data, &pdf.Share)
		if err != nil {
			return nil, err
		}

		pdfList = append(pdfList, pdf)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pdfList, nil
}
