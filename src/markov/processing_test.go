package markov

import (
	"os"
	"testing"

	"github.com/crunchypi/markov-go-sql.git/src/storage"
	"github.com/crunchypi/markov-go-sql.git/src/storage/neo4j"
	"github.com/crunchypi/markov-go-sql.git/src/storage/sqlite"
)

const (
	datapath = "testdata/test.txt"
	dbpath   = "testdata/test.sqlite"
)

var currentTestDB = newDBNeo4j

func newDBSQLite() storage.DBAbstracter {
	db, err := sqlite.New(dbpath)
	if err != nil {
		panic("db preparation failed")
	}
	return db
}

func newDBNeo4j() storage.DBAbstracter {
	const (
		uri = "bolt://localhost:7687"
		usr = ""
		pwd = ""
		enc = false
	)

	db, err := neo4j.New(uri, usr, pwd, enc)
	if err != nil {
		panic("failed while opening neo4jDB. credentials?")
	}
	return db
}

func cleanup() {
	os.Remove(dbpath)
}

func TestNew(t *testing.T) {
	_, err := New(datapath, currentTestDB())
	defer cleanup()
	if err != nil {
		t.Error("failed while creating new MarkovChain obj: ", err)
	}
}

func TestProcessCorpusCrashTest(t *testing.T) {
	m, _ := New(datapath, currentTestDB())
	defer cleanup()
	// # check for crash
	for i := -5; i < len(m.corpus)+10; i++ {
		m.ProcessCorpusByOrder(i, false)
	}
}

// Test is verified by checking db
func TestProcessCorpusByOrder(t *testing.T) {
	cleanup()
	m, _ := New(datapath, currentTestDB())
	// m.ProcessCorpusByOrder(-1, false)
	m.ProcessCorpusComplete(false)

}
