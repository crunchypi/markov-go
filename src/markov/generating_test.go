package markov

import (
	"testing"

	"github.com/crunchypi/markov-go-sql.git/src/protocols"
	"github.com/crunchypi/markov-go-sql.git/src/storage/sqlite"
)

// ----------------------------------------------------------------
//
// # This test deals with probability - built for manual inspection.
//
// ----------------------------------------------------------------

const (
	datapathGen = "testdata/test.txt"
	dbpathGen   = "testdata/genTest.sqlite"
)

func TestChoose(t *testing.T) {
	words := []string{}
	scores := []int{}

	if len(words) != len(scores) {
		t.Error("test setup issue")
	}

	relShips := make([]protocols.WordRelationship, 0)
	for i := 0; i < len(words); i++ {
		relShips = append(relShips, protocols.WordRelationship{
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
	db, err := sqlite.New(dbpathGen)
	if err != nil {
		t.Error("failed while creating db object:", err)
	}
	m, err := New(datapathGen, db)
	if err != nil {
		t.Error("faild while setting up MChain:", err)
	}

	m.ProcessCorpusComplete(false)
	res := m.GenerateSimple("some", 2)
	t.Log("#Result:", res)

}
