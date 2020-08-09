package markov

import (
	"math/rand"
	"time"

	"github.com/crunchypi/markov-go-sql.git/src/storage"
)

// GenerateSimple uses a seed to generate random sentences with length of 'wordcount',
// using the database connection set in the struct. Negative 'wordCount' will keep it
// going until a stored word doesn't have an outgoing edge.
// Methodology is a random weighted walk.
func (m *MarkovChain) GenerateSimple(seed string, wordCount int) []string {
	// # not setting cap in make func because negative wordCount is allowed.
	result := make([]string, 0)
	result = append(result, seed)

	i := 0
	for {
		// # Stop condition no. 1. Not Normal i;bool;i++ form because
		// # the limit can be disabled with negative wordCount.
		if i >= wordCount-1 && wordCount > 0 {
			break
		}
		// # Disable iteration stop limit if wordCount is negative.
		if i >= 0 {
			i++
		}
		// # Use last item in chain to fetch candidates to choose from.
		collection, _ := m.db.SucceedingX(result[len(result)-1])
		choice, ok := choose(collection)
		// # Stop condition no. 2: If chain is broken (no more succeeding)
		if !ok {
			return result
		}
		// # Setup for next pick.
		result = append(result, choice)
	}
	return result
}

// # Weighted choice.
func choose(choices []storage.WordRelationship) (string, bool) {
	// # Only fail condition: Algorithm needs len to be > 0
	// # because there must be something to choose from.
	if len(choices) == 0 {
		return "", false
	}

	// # Set distribution of probabilities.
	selectionPool := make(map[int]string, len(choices))
	max := 0
	for i := 0; i < len(choices); i++ {
		newMax := max + choices[i].Count
		selectionPool[newMax] = choices[i].Other
		max = newMax
	}

	// # pick.
	rand.Seed(time.Now().UnixNano())
	selectedInt := rand.Intn(max)

	for key, val := range selectionPool {
		if key >= selectedInt {
			return val, true
		}
	}

	// # This line is required but not expected to run,
	// # as the only known fail condition is handled at
	// # the start of this function.
	panic("unexpected end of function")
}
