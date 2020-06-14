package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// import mt "crunchypi/markov-tools/other"

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
	if len(args) == 3 && args[0] == "-train" {
		subSlice := args[1:]
		from, okFrom := parseArgs(subSlice, "-from")
		to, okTo := parseArgs(subSlice, "-to")
		if okFrom && okTo {
			process(from, to)
			fmt.Println("\n\nDone.")
			return
		}
	}

	// # Try Retrieving data.
	if len(args) == 4 && args[0] == "-chat" {
		subSlice := args[1:]
		from, okFrom := parseArgs(subSlice, "-from")
		init, okInit := parseArgs(subSlice, "-init")
		n, okLen := parseArgs(subSlice, "-len")
		if okFrom && okInit && okLen {
			n = "3"
			nVal, err := strconv.ParseInt(n, 0, 64)
			if err == nil {
				fmt.Println(retrieve(from, init, int(nVal)))
				return
			}
		}
	}

	// # Default help.
	fmt.Println(
		"Incorrect arguments. Example:\n",
		"<call program> -train -from<txtfile> -to=<dbfile> 				  ### Train\n",
		"<call program> -chat  -from<dbfile> -init=<firstword> -len=<int> ### Generate",
	)

}
