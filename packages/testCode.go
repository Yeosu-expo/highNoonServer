package packages

import "log"

func printChunk(chunk PlayerChunk) {
	log.Println("-------Player Chunk-------")
	log.Println("ID:", chunk.Id)
	log.Println("Time:", chunk.DateTime)
	log.Println("Score:", chunk.Score)
	log.Println("--------------------------")
}
