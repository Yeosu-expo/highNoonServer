package packages

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func ServingChunkHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()

	var chunk PlayerChunk
	if err = json.Unmarshal(data, &chunk); err != nil {
		log.Println(err)
		return
	}

	db, err := sql.Open("mysql", "root:9250@tcp(localhost:3306)/highnoon")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	query := `INSERT INTO PlayerChunk VALUES (?,?,?)`
	_, err = db.Exec(query, chunk.Id, chunk.DateTime, chunk.Score)
	if err != nil {
		log.Println("Failed to insert row:", err)
		return
	}

	log.Println("Player Chunk Recorded.")
	printChunk(chunk)
}
