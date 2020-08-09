package markov

import (
	"io/ioutil"
	"strings"

	"github.com/crunchypi/markov-go-sql.git/src/storage"
)

type MarkovChain struct {
	corpus []string
	db     storage.DBAbstracter
}

func New(db storage.DBAbstracter) *MarkovChain {
	mc := MarkovChain{}
	mc.db = db
	return &mc
}

func (m *MarkovChain) ReadFileContent(path string) error {
	content, err := ioutil.ReadFile(path)
	m.corpus = strings.Split(string(content), " ")
	return err
}
