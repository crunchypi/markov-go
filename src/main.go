package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/crunchypi/markov-go-sql.git/src/markov"
)

func main() {
	CLI()
}

func parseArgs(args []string, search string) (string, bool) {
	// # Look through arguments.
	for _, arg := range args {
		// # Separate chunk into identifier=target.
		dual := strings.Split(arg, "=")
		if dual[0] == search {
			// # Return target and success.
			return dual[1], true
		}
	}
	return "", false
}

// CLI : Command Line Interface.
func CLI() {
	args := os.Args[1:]

	// # Try training DB.
	if len(args) == 4 && args[0] == "-train" {
		subSlice := args[1:]
		from, okFrom := parseArgs(subSlice, "-from")
		to, okTo := parseArgs(subSlice, "-to")
		order, okOrd := parseArgs(subSlice, "-order")
		if okFrom && okTo && okOrd {
			orderVal, err := strconv.ParseInt(order, 0, 64)
			if err == nil {
				markov.Process(from, to, int(orderVal))
				fmt.Println("\n\nDone.")
				return
			}
		}
	}

	// # Try Retrieving data.
	if len(args) == 4 && args[0] == "-chat" {
		subSlice := args[1:]
		from, okFrom := parseArgs(subSlice, "-from")
		init, okInit := parseArgs(subSlice, "-init")
		n, okLen := parseArgs(subSlice, "-len")
		if okFrom && okInit && okLen {
			nVal, err := strconv.ParseInt(n, 0, 64)
			if err == nil {
				fmt.Println(markov.Retrieve(from, init, int(nVal)))
				return
			}
		}
	}

	// # Default help.
	fmt.Println(
		"Incorrect arguments. Example:\n",
		"<call program> -train -from=<txtfile> -to=<dbfile> -order=<int>   ### Train\n",
		"<call program> -chat  -from=<dbfile> -init=<firstword> -len=<int> ### Generate",
	)

}
