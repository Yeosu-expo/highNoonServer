package packages

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const (
	DB = "root:9250@tcp(localhost:3306)/"
)

func getDB(DBName string) *sql.DB {
	db, err := sql.Open("mysql", DBName)
	if err != nil {
		log.Println(err)
		return nil
	}
	return db
}

func HashingPassword(password string) (string, error) {
	hased, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return "", err
	}
	return string(hased), err
}
