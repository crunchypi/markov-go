package markov

import (
	"math/rand"
	"time"

	"github.com/crunchypi/markov-go-sql.git/src/protocols"
)

func (m *MarkovChain) GenerateSimple(seed string, wordCount int) []string {
	// # not setting cap in make func because negative wordCount is allowed.
	result := make([]string, 0)
	result = append(result, seed)

	i := 0
	for {
		// # Stop condition. Not Normal i;bool;i++ form because
		// # the limit can be disabled with negative wordCount.
		if i >= wordCount-1 && wordCount > 0 {
			break
		}
		// # Disable iteration stop limit if wordCount is negative.
		if i >= 0 {
			i++
		}
		// # Use last item in chain to fetch candidates to choose from.
		collection := m.db.SucceedingX(result[len(result)-1])
		choice, ok := choose(collection)
		// # stop if chain is broken (no more succeeding)
		if !ok {
			return result
		}
		// # Setup for next pick.
		result = append(result, choice)
	}
	return result
}

func choose(choices []protocols.WordRelationship) (string, bool) {
	// # Only fail condition: Algorithm needs len to be > 0
	// # because there must be something to choose from.
	if len(choices) == 0 {
		return "", false
	}

	// # Each word has a range associated to it.
	// # Example with two words:
	// # 	'first' has score of 10
	// # 	'second' has score of 20.
	// #
	// #	.. So map will be map{
	// #				'first':10,
	// #				'second':30,
	// #			}
	// # 	.. which means that 'first' has a range 0-10
	// # 	and 'second' has a range of 10-30. note that
	// # 	order isn't important. Example continued below.
	// #
	selectionPool := make(map[int]string, len(choices))
	max := 0
	for i := 0; i < len(choices); i++ {
		newMax := max + choices[i].Count
		selectionPool[newMax] = choices[i].Other
		max = newMax
	}

	// # Continuing example:
	// # Use the highest value in the map (30)
	// # in the random Intn function. Probability
	// # is distributed proportionally, so 'second'
	// # is 2x time likely to be chosen than 'first'.
	rand.Seed(time.Now().UnixNano())
	selectedInt := rand.Intn(max)

	for key, val := range selectionPool {
		if key >= selectedInt {
			return val, true
		}
	}

	panic("unexpected end of function")
}
