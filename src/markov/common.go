package markov

import (
	"io/ioutil"
	"strings"

	"github.com/crunchypi/markov-go-sql.git/src/protocols"
)

var _ protocols.MarkovAbstracter = (*MarkovChain)(nil)

type MarkovChain struct {
	corpus []string
	db     protocols.DBAbstracter
}

func New(dataPath string, db protocols.DBAbstracter) (*MarkovChain, error) {
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
