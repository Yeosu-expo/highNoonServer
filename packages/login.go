package packages

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()

	var chunk UserInfoChunk
	if err = json.Unmarshal(data, &chunk); err != nil {
		log.Println(err)
		return
	}

	log.Println(chunk)
	db := getDB(DB + "highnoon")
	if db == nil {
		return
	}
	defer db.Close()

	hased, err := HashingPassword(chunk.User_Password)
	if err != nil {
		log.Println(err)
		return
	}

	query := `INSERT INTO logintable VALUES (?,?)`
	_, err = db.Exec(query, chunk.User_ID, hased)
	if err != nil {
		log.Println("Failed to insert row:", err)
		return
	}

	if err = mkMMRRow(db, chunk.User_ID); err != nil {
		log.Println(err)
		return
	}

	log.Println(chunk.User_ID, "Sign Up.")
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()

	var chunk UserInfoChunk
	if err = json.Unmarshal(data, &chunk); err != nil {
		log.Println(err)
		return
	}

	db := getDB(DB + "highnoon")
	if db == nil {
		return
	}
	defer db.Close()

	query := `SELECT user_password from logintable WHERE id = ?`
	rows, err := db.Query(query, chunk.User_ID)
	if err != nil {
		log.Println("Failed to insert row:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user_password string
		if err := rows.Scan(&user_password); err != nil {
			log.Println(err)
			return
		}
		log.Println(user_password)

		res := bcrypt.CompareHashAndPassword([]byte(user_password), []byte(chunk.User_Password))
		if res == nil {
			w.WriteHeader(http.StatusOK)
			log.Println(chunk.User_ID, "Sign In.")
			return
		}
	}

	log.Println(chunk.User_ID, "Sign In Failed.")
	w.WriteHeader(http.StatusNotFound)
	return
}

func mkMMRRow(db *sql.DB, uid string) error {
	query := "INSERT INTO mmr VALUES(?, 100, 0, 0)"
	_, err := db.Exec(query, uid)
	return err
}
