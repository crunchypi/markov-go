package main

import (
	"log"

	"github.com/crunchypi/markov-go-sql.git/src/markov"
	"github.com/crunchypi/markov-go-sql.git/src/storage/neo4j"
)

func main() {
	const (
		uri = "bolt://localhost:7687"
		usr = ""
		pwd = ""
		enc = false
	)

	db, err := neo4j.New(uri, usr, pwd, enc)
	if err != nil {
		log.Fatal(err)
	}

	mc, err := markov.New("../data/a.txt", db)
	if err != nil {
		log.Fatal(err)
	}

	mc.ProcessCorpusByOrder(2, true)
}
