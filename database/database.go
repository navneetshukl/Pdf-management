package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/joho/godotenv"
)

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

func InsertUser(name,email,password string){
	conn,err:=DB_Connect()
	defer conn.Close()
	if err!=nil{
		panic(err)
	}
	query:=`insert into users(name,email,password) values($1,$2,$3)`

	_,err =conn.Exec(query,name,email,password)

	if err!=nil{
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
