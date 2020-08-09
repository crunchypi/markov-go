package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/crunchypi/markov-go-sql.git/src/storage"

	"github.com/crunchypi/markov-go-sql.git/src/markov"
	"github.com/crunchypi/markov-go-sql.git/src/storage/neo4j"
	"github.com/crunchypi/markov-go-sql.git/src/storage/sqlite"
)

// # See init for var details.
var (
	runType  string
	dataPath string
	order    int
	dbChoice string

	n4juri string
	n4jusr string
	n4jpwd string

	sqlitePath string

	seed      string
	wordCount int
)

// # Used to avoid default value of int.
var invalidDefaultInt = -789

// # Set all flag vars.
func init() {
	// # process .txt document or generate txt to stdout
	flag.StringVar(&runType, "run", "",
		": 'train' creates a markov chain. \n"+
			": 'generate' uses a markov chain to generate text.")

	// # path to .txt document containing txt to process
	flag.StringVar(&dataPath, "datapath", "",
		": Choose a path to a text file (chain training).")
	// # order for processing in markov chain
	flag.IntVar(&order, "order", invalidDefaultInt,
		": Traditional order of a markov chain (chain training).")

	// # How to store processed word relationships.
	flag.StringVar(&dbChoice, "db", "",
		": Choose db storage. Choices are sqlite and neo4j\n"+
			"  (neo4j requires credentials and must be running)")

	// # credentials for neo4j
	flag.StringVar(&n4juri, "uri", "", ": URI for neo4j.")
	flag.StringVar(&n4jusr, "usr", "", ": Username for neo4j.")
	flag.StringVar(&n4jpwd, "pwd", "", ": Password for neo4j.")

	// # path to sqlite
	flag.StringVar(&sqlitePath, "sqlite", "",
		": Path to sqlite file.")

	flag.StringVar(&seed, "seed", "",
		": Initial word when using '-run generate'")
	flag.IntVar(&wordCount, "wordcount", invalidDefaultInt,
		": Dictates how many words '-run generate' will make")

	flag.Parse()

}

// # quick exit.
func exit(msg string) {
	help(msg)
	os.Exit(1)
}

// # print out help
func help(msg string) {
	// # Specific messages.
	line := "\n"
	for i := 0; i < 120; i++ {
		line += "-"
	}
	line += "\n"
	fmt.Println(line, msg, line)
	// # Flag descriptions.
	flag.PrintDefaults()
	// # Examples
	lines := []string{
		"\n----EXAMPLES----",

		"\nCreate a markov chain with sqlite:",
		"\n\t-run train -db sqlite -sqlite ./s.sqlite -datapath ./x.txt -order 3",

		"\nGenerate text with sqlite:",
		"\n\t-run generate -db sqlite -sqlite ./s.sqlite -wordcount 4 -seed help",

		"\nCreate a markov chain with neo4j:",
		"\n\t-run train -db neo4j -uri bolt://localhost:7687 " +
			"-usr neo4j, -pwd neo4j -datapath ./x.txt -order 3",

		"\nGenerate text with neo4j:",
		"\n\t-run generate -db neo4j -uri bolt://localhost:7687 " +
			"-usr neoj4, -pwd neo4j -wordcount 6 -seed me",

		"\n\n\n",
	}
	fmt.Println(lines)
}

// # Helps with formatting specific error messages. Usually used
// # as an argument to help().
func msg(flag string, flagVar interface{}, options ...string) string {
	opt := ``
	for _, v := range options {
		opt += v + ","
	}
	opt = opt[0 : len(opt)-1] // # remove last comma.
	return fmt.Sprintf(
		"Invalid input for -%v: '%v'. Options: [%v]",
		flag, flagVar, opt,
	)
}

// Run is the point of entry for the CLI.
func Run() {
	switch runType {
	case "train":
		startProcessing()
	case "generate":
		generate()
	default:
		help(msg("runType", runType, "train, generate"))
	}
}

// # used with a -run=generate
func generate() {
	m := markov.New(chooseDB())
	if wordCount == invalidDefaultInt {
		exit(msg("wordcount", wordCount,
			fmt.Sprintf("Any number except for default: %v",
				invalidDefaultInt)))
	}
	if seed == "" {
		exit(msg("seed", seed, "Any words"))
	}

	fmt.Println(m.GenerateSimple(seed, wordCount))
}

// # used with a -run=train flag
func startProcessing() {
	m := markov.New(chooseDB())
	err := m.ReadFileContent(dataPath)
	if err != nil {
		exit(
			msg("datapath", dataPath, "rel/abs path") +
				"\n\tError while processing: " + err.Error(),
		)
	}
	if order == invalidDefaultInt {
		exit(msg("order", order,
			fmt.Sprintf("Any number except for default: %v",
				invalidDefaultInt)))
	}
	m.ProcessCorpusByOrder(order, true)

}

// # used to handle -db flag
func chooseDB() storage.DBAbstracter {

	switch dbChoice {
	case "sqlite":
		return prepSQLite()
	case "neo4j":
		return prepNeo4j()
	}
	exit(msg("db", dbChoice, "sqlite", "neo4j"))
	return nil
}

// # used when -db=neo4j
func prepNeo4j() storage.DBAbstracter {
	if n4juri == "" {
		exit(msg("uri", n4juri, "Any uri. Usually bolt://localhost:7687"))
	}
	if n4jusr == "" {
		exit(msg("usr", n4jusr, "Any username. Default is neo4j"))
	}
	if n4jpwd == "" {
		exit(msg("pwd", n4jpwd, "Any pwd. Default is neo4j"))
	}

	db, err := neo4j.New(n4juri, n4jusr, n4jpwd, false)
	if err != nil {
		exit(fmt.Sprintf("Error while setting up Neo4j: %s", err))
	}
	fmt.Println("WARN: make sure credentials are correct and server is running.")
	return db
}

// # used when -db=sqlite
func prepSQLite() storage.DBAbstracter {
	if sqlitePath == "" {
		exit(msg(
			"sqlite",
			sqlitePath,
			"Any path to an sqlite file made with this tool",
		))
	}
	db, err := sqlite.New(sqlitePath)
	if err != nil {
		exit(msg("-sqlite", sqlitePath, "relative/abs path"))
	}
	return db
}

func choiceGen() {

}
