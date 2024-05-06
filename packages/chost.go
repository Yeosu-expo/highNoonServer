package packages

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func InsertChostHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()

	var chost ChostData
	if err = json.Unmarshal(data, &chost); err != nil {
		log.Println(err)
		return
	}

	if err = insertChostData(chost); err != nil {
		log.Println(err)
		return
	}

	log.Println(chost.UID, "chost data inserted.")
}

func GetChostHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	chost := getChostData(uid)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(chost)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("--------GET CHOST-------")
	log.Println(chost)
	log.Println("------------------------")
}

func insertChostData(chost ChostData) error {
	db := getDB(DB + "highnoon")
	defer db.Close()

	var formtype FormType
	formtype.Form = chost.Form
	var timingtype TimingType
	timingtype.Timing = chost.Timing

	form, err := json.Marshal(formtype)
	if err != nil {
		return err
	}
	timing, err := json.Marshal(timingtype)
	if err != nil {
		return err
	}

	query := "insert into chost values(?, ?, ?, ?, ?)"
	if _, err = db.Exec(query, chost.UID, chost.Accuracy, form, timing, chost.Index); err != nil {
		return err
	}

	return nil
}

func getChostData(uid string) ChostData {
	db := getDB(DB + "highnoon")
	defer db.Close()

	query := "select * from chost where uid=(?)"
	rows, err := db.Query(query, uid)
	if err != nil {
		log.Println(err)
		return ChostData{}
	}

	var data []ChostData
	for rows.Next() {
		var chost ChostData
		var form []byte
		var timing []byte
		if err := rows.Scan(&chost.UID, &chost.Accuracy, &form, &timing, &chost.Index); err != nil {
			log.Println(err)
			return ChostData{}
		}

		var formtype FormType
		if err := json.Unmarshal(form, &formtype); err != nil {
			log.Println(err)
			return ChostData{}
		}

		var timingtype TimingType
		if err := json.Unmarshal(timing, &timingtype); err != nil {
			log.Println(err)
			return ChostData{}
		}

		chost.Form = formtype.Form
		chost.Timing = timingtype.Timing

		data = append(data, chost)
	}

	index := rand.Intn(len(data))

	return data[index]
}
