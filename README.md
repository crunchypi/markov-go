# markov-go-sql
Markov chain in go - storage with SQLite3


Requirements:
	
	Go 	(created with 1.14)
	SQLite3
	github.com/mattn/go-sqlite3


Usage:

	To create the markov chain DB:
		<call program> -train -from=<path_to_file.txt> -to=<path_to_file.db -order=<int>>
	
	Note: 	textfile/datasource needs to exist, db file is created if it's missing.
		'order' is the n-gram order, by convention, at the moment it should be > 2.

	To use data from a DB file:
		<call program> -chat -from<path_to_file.db> -init=<first_word> -len=<int>

	Note: 	'init' means the initial state of the chain, as in; the program will use that
		word as an inspiration. 'len' means how many sequantial words the program should
		fetch.
	

DB IO is highly inefficient at the moment, the status of this project is 'prototype'. 
