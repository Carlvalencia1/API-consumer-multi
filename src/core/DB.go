package core

import (
	"database/sql"
	"log"
	"os"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Cant load the .env!")
	}
}

func NewMysql() *sql.DB {
	LoadEnv()

	//env variables;
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Error to connect to the Mysql1: %v", err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Succes connection to database!")
	return db
}