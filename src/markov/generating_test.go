package markov

import (
	"testing"

	"github.com/crunchypi/markov-go-sql.git/src/storage"
)

// ----------------------------------------------------------------

// # This test deals with probability - built for manual inspection.

// ----------------------------------------------------------------

const (
	datapathGen = "testdata/test.txt"
	dbpathGen   = "testdata/genTest.sqlite"
)

// # Choices : newDBSQLite, newDBNeo4j. Functions
// # defined in processing_test.go
var currentTestDB_ = newDBNeo4j

func TestChoose(t *testing.T) {
	words := []string{}
	scores := []int{}

	if len(words) != len(scores) {
		t.Error("test setup issue")
	}

	relShips := make([]storage.WordRelationship, 0)
	for i := 0; i < len(words); i++ {
		relShips = append(relShips, storage.WordRelationship{
			Other: words[i],
			Count: scores[i],
		})
	}

	choice, ok := choose(relShips)
	t.Log("## CHOSE:", choice, ok)
	if ok {
		t.Error("should not be ok")
	}
}

func TestGenerateSimple(t *testing.T) {

	m, err := New(datapathGen, currentTestDB_())
	if err != nil {
		t.Error("faild while setting up MChain:", err)
	}

	m.ProcessCorpusByOrder(-1, true)
	res := m.GenerateSimple("some", 4)
	t.Log("#Result:", res)

}
