package markov

import (
	"io/ioutil"
	"strings"

	"github.com/crunchypi/markov-go-sql.git/src/storage"
)

// MarkovChain holds a db connection which is used to either load
// a markov chain structure into storage, or generate words from
// that storage.
type MarkovChain struct {
	corpus []string
	db     storage.DBAbstracter
}

// New returns a ref to a MarkovChain obj, db handler is required.
func New(db storage.DBAbstracter) *MarkovChain {
	mc := MarkovChain{}
	mc.db = db
	return &mc
}

// ReadFileContent is a simple method which reads a text file
// into the corpus attribute.
func (m *MarkovChain) ReadFileContent(path string) error {
	content, err := ioutil.ReadFile(path)
	m.corpus = strings.Split(string(content), " ")
	return err
}
