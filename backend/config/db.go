package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() *sql.DB {
	connectionString := "root:password@tcp(localhost:3000)/mymariadb?parseTime=true"
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	return db
}
