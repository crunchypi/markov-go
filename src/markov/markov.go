package markov

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/crunchypi/markov-go-sql.git/src/dbport"
)

// Retrieve uses the markov chain technique to retrieve words.
// 	- DBName 	: Name of the DB (generate with Process())
//	- initial	: First string to kick-off the returned sequence.
//	- n 		: How many words to retrieve (not including initial)
func Retrieve(DBName string, initial string, n int) []string {

	result := make([]string, 0, n)   // # Store result.
	result = append(result, initial) // # Add initial str.

	last := initial // # Tracking
	for i := 0; i < n; i++ {
		// # Fetch results from DB.
		others, nums := dbport.SucceedingX(DBName, last)
		// # Choose a weighted random - do nothing if result is empty.
		choice, ok := weightedChoice(others, nums)
		if ok {
			result = append(result, choice)
			last = choice
		}
	}
	return result
}

// Process populates an SQLite DB with content from a
// specified file in a markov-chain fashion.
// 	- fnText 	: input filename (should be .txt)
//	- DBName 	: output filename (.db/.sqlite)
//  - order		: amount of related words. (NOTE: < 3 causes issues ATM).
func Process(fnText string, DBName string, order int) {

	// # Slice of words -> slice of slices (windows)
	words := strings.Split(readFileContent(fnText), " ")
	windows := ngram(words, order)
	if len(windows) == 0 {
		fmt.Println("Not enough word combinations, aborting.")
		return
	}

	// # Prepare DB. Won't overwrite if exists.
	dbport.CreateDictionary(DBName)

	for i, window := range windows {
		// # Progress feedback
		fmt.Printf("\r Processing. Chunks remaining: %v", len(windows)-i)

		current, others := window[0], window[1:]

		// # Count each relationship between current and others.
		for _, other := range others {
			// @ score eq: float32(windowSize - j - 1)
			switch dbport.CheckWordPair(DBName, current, other) {
			case true:
				dbport.IncrementWordPair(DBName, current, other)
			case false:
				dbport.InsertNewWord(DBName, current, other, 1)
			} // # case.
		} // # others loop.
	} // # window loop.
}

// Creates 'windows of words'. Example:
// 		content	: ["some", "random", "string", "content"]
//		order	: 3
//		result	: [[some random string] [random string content]]
//
// NOTE: @ Bug on order: < 3
func ngram(content []string, order int) [][]string {
	result := make([][]string, 0, len(content)/order)
	adjacentCount := (order - 1) / 2
	for i := adjacentCount; i < len(content)-adjacentCount; i++ {
		window := content[i-adjacentCount : i+adjacentCount+1]
		result = append(result, window)
	}
	return result
}

// Plainly reads a file. Returns empty str on fail.
func readFileContent(name string) string {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(content)
}

// Choose a random weighted value. Example:
//		vals	: ["random", "words"]
//		nums	: [2,1]
//		result	: probability of "random" is 33...%
//
// NOTE 1: length of 'vals' and 'nums' must be the same.
// NOTE 2: algo is highly inefficient.
func weightedChoice(vals []string, nums []float32) (string, bool) {
	rand.Seed(time.Now().UnixNano())
	// Ensure symmetry.
	if len(vals) != len(nums) && len(vals) < 1 {
		return "", false
	}

	// # Create a pool of words such that the occurance
	// # of each word = its corresponding position in nums.
	pool := make([]string, 0, len(vals)*10) // len(vals) * 10 is arbitrary.
	for i, v := range vals {
		for j := 0; j < int(nums[i]); j++ {
			pool = append(pool, v)
		}
	}

	// # Choose.
	return pool[rand.Intn(len(pool))], true
}
