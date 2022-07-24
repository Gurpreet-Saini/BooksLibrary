package driver

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectToSQL() (*sql.DB, error) {
	db, err := sql.Open("mysql",
		"root:Gurpreet@0848@tcp(localhost:3306)/test")
	if err != nil {
		log.Println(err)

		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Println(err)

		return nil, err
	}

	log.Println("Connected")

	return db, nil
}
