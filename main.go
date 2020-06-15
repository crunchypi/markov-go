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
	// fn := "db/col.db"
	// tableCreate(fn, "dank")
	// res := tableExists(fn, "dank")
	// fmt.Println(res)
	// pairInsert(fn, "dank", "bank", 1)
	// res := pairExists(fn, "dank", "bank", true)
	// fmt.Println(res)
	// pairIncrement(fn, "dank", "bank", 1)
	// n, v := tableDump(fn, "stank")
	// fmt.Println(n, v)
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
				process(from, to, int(orderVal))
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
				fmt.Println(retrieve(from, init, int(nVal)))
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
