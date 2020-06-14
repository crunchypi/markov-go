package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func retrieve(fnDB string, initial string, n int) []string {
	// # Accumulates relevant words.
	result := make([]string, 0, n)
	last := initial

	for i := 0; i < n; i++ {
		rand.Seed(time.Now().UnixNano())
		// @ TODO implement weighted choice.
		others, _ := succeedingX(fnDB, last)
		// # Random choice from others is added to result.
		length := len(others)
		if length > 0 {
			r := rand.Intn(length-0) + 0
			choice := others[r]
			result = append(result, choice)
			last = choice
		}
	}
	return result
}

func process(fnText, fnDB string) {
	// # Get long string from file.
	content := readFileContent(fnText)
	// # Slice of words.
	words := strings.Split(content, " ")
	// # Slice of slices (moving window).
	windows := ngram(words, 4)
	// # DB created but, does not delete old table if exists.
	createDictionary(fnDB)
	// # Counter for printout.
	counter := 0
	// # Abort on empty windows.
	if len(windows) == 0 {
		fmt.Println("Not enough word combinations, aborting.")
		return
	}
	for _, window := range windows {
		fmt.Printf("\r Num of processed chunks: %v", counter)
		current := window[0]
		others := window[1:]
		// # Count each relationship between current and others.
		for _, other := range others {
			// @ score eq: float32(windowSize - j - 1)
			exists := checkWordPair(fnDB, current, other)
			if exists {
				incrementWordPair(fnDB, current, other)
			} else {
				insertNewWord(fnDB, current, other, 1)
			}
		}
		counter++
	}
}

func ngram(content []string, order int) [][]string {
	// @ Bug on order: < 3
	result := make([][]string, 0, len(content)/order)
	adjacentCount := (order - 1) / 2
	for i := adjacentCount; i < len(content)-adjacentCount; i++ {
		window := content[i-adjacentCount : i+adjacentCount+1]
		result = append(result, window)
	}
	return result
}

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
