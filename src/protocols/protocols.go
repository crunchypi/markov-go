package protocols

// DBAbstracter abstracts implementations of any DB.
// Purpose is to let a MarkovAbstracter save and retrieve
// word relationship data.
type DBAbstracter interface {
	IncrementPair(word, other string, dst int)
	// PairExists()
	SucceedingX(word string) []WordRelationship
}

// MarkovAbstracter abstracts implementation of markov
// chain implementation.
type MarkovAbstracter interface {
	ProcessCorpusByOrder(order int, verbose bool)
	ProcessCorpusComplete(verbose bool)
	GenerateSimple(seed string, wordCount int) []string
}

type WordRelationship struct {
	Word     string
	Other    string
	Distance int
	Count    int
}
