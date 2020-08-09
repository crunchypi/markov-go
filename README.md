# markov-go-sql
Markov chain in go - storage with SQLite3/Neo4j


Requirements:
	Go 	(created with 1.14)
	
	If using sqlite:
		SQLite3 on machine (created with 3.28.0)
		github.com/mattn/go-sqlite3
		
	If using neo4j:
		Neo4j on machine (created with 4.1.1)
		https://github.com/neo4j/neo4j-go-driver


Usage CLI (Load into db):

	// Create a markov chain with sqlite:
		-run train -db sqlite -sqlite ./s.sqlite -datapath ./somefile.txt -order 3
		
	// Create a markov chain with neo4j:
		-run train -db neo4j -uri bolt://localhost:7687 -usr <username> -pwd <password> -datapath ./somefile.tx -order 3"
		
Usage CLI (Generate from db):

	// Generate text with sqlite
		-run generate -db sqlite -sqlite ./s.sqlite -wordcount 4 -seed help

	// Generate text with neo4j:
		-run generate -db neo4j -uri bolt://localhost:7687 -usr <username> -pwd <password> -wordcount 6 -seed help



Usage API (Load into sqlite):
	
	// prepare db:
	sqlitePath := "./whereToStoreSQLiteFile.sqlite"
	db := sqlite.New(sqlitePath) // implements storage.DBAbstracter
	
	// prepare markov chain processor:
	txtPath := "./whereTxtFileIs.txt"
	mc := markov.New(db) // takes in storage.DBAbstracter
	mc.ReadFileContent(txtPath)
	
	// process:
	mc.ProcessCorpusByOrder(3, true) // 3 is order, true is verbosity.
	
	// Alternative process, order is not required:
	mc.ProcessCorpusComplete(true) // true is verbosity.
	
	
Usage API (Load into Neo4j):

	// prepare db:
	uri, usr, pwd, enc := "bolt://localhost:7687", "neo4j", "neo4j", false
	db := neo4j.New(uri, usr, pwd, enc) // implements storage.DBAbstracter

	// prepare markov chain processor:
	txtPath := "./whereTxtFileIs.txt"
	mc := markov.New(db) // takes in storage.DBAbstracter
	mc.ReadFileContent(txtPath)
	
	// process:
	mc.ProcessCorpusByOrder(3, true) // 3 is order, true is verbosity.
	
	// Alternative process, order is not required:
	mc.ProcessCorpusComplete(true) // true is verbosity.
	
	
Usage API (generate from sqlite):
	
	// prepare db:
	sqlitePath := "./whereSQLiteFileIsLocated.sqlite"
	db := sqlite.New(sqlitePath) // implements storage.DBAbstracter

	// prepare markov chain generator:
	mc := markov.New(db) // takes in storage.DBAbstracter
	
	seed, wordCount := "help", 3
	result := mc.GenerateSimple(seed, wordCount)
	
	
Usage API (generate from neo4j):

	// prepare db:
	uri, usr, pwd, enc := "bolt://localhost:7687", "neo4j", "neo4j", false
	db := neo4j.New(uri, usr, pwd, enc) // implements storage.DBAbstracter
	
	// prepare markov chain generator:
	mc := markov.New(db) // takes in storage.DBAbstracter

	seed, wordCount := "help", 3
	result := mc.GenerateSimple(seed, wordCount)
	
	
	
	

	
	
	
