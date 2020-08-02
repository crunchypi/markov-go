package storage

// DBAbstracter abstracts implementations of any DB.
// Purpose is to let a MarkovAbstracter save and retrieve
// word relationship data.
type DBAbstracter interface {
	IncrementPair(word, other string, dst int) error
	SucceedingX(word string) ([]WordRelationship, error)
}

// WordRelationship stores relationships of words.
// - 'Word' and 'Other' are the words
// - 'Distance' describes how far the words are from
//   eachother (in a sentence, for example).
// - 'Count' indicates how many such relationships
//   there are.
type WordRelationship struct {
	Word     string
	Other    string
	Distance int
	Count    int
}
