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

func New(dataPath string, db storage.DBAbstracter) (*MarkovChain, error) {
	mc := MarkovChain{}
	mc.db = db

	// # Try load corpus.
	err := mc.readFileContent(dataPath)
	if err != nil {
		return &mc, err
	}

	return &mc, nil
}

func (m *MarkovChain) readFileContent(path string) error {
	content, err := ioutil.ReadFile(path)
	m.corpus = strings.Split(string(content), " ")
	return err
}
