package markov

import (
	"os"
	"strings"
	"testing"

	"github.com/crunchypi/markov-go-sql.git/src/storage/sqlite"
)

const (
	datapath = "testdata/test.txt"
	dbpath   = "testdata/test.sqlite"
)

var dataContent = "some random string content\n" // # hast to be the same data in test.txt

var currentTestDB = newDBSQLite

func newDBSQLite() *sqlite.SQLiteManager {
	db, err := sqlite.New(dbpath)
	if err != nil {
		panic("db preparation failed")
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

func TestNewData(t *testing.T) {
	m, _ := New(datapath, currentTestDB())
	defer cleanup()

	c := strings.Split(dataContent, " ")
	for i := 0; i < len(c); i++ {
		if c[i] != m.corpus[i] {
			t.Error("data inconsistency on iter:", i)
		}
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
	m.ProcessCorpusByOrder(2, false)
}

// Test is verified by checking db
func TestProcessCorpusComplete(t *testing.T) {
	cleanup()
	m, _ := New(datapath, currentTestDB())

	m.ProcessCorpusComplete(true)

}
