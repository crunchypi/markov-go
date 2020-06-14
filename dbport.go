package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// func TestDB(t *testing.T) {

// 	name := "sqlite.db"
// 	createDB(name, true)
// 	createDictionary(name)
// 	insertNewWord(name, "dank", "bank", 32.0)
// 	incrementWordPair(name, "dank", "bank")
// 	fmt.Println(checkWordPair(name, "dak", "bank"))

// }

func createDB(DBName string, overwriteOld bool) bool {
	// # Remove old DB file.
	if overwriteOld {
		os.Remove(DBName)
		file, err := os.Create(DBName)
		if err != nil {
			log.Fatal(err.Error())
			return false
		}
		file.Close()
	}

	// # Create new DB file.
	sqliteDB, _ := sql.Open("sqlite3", DBName)
	defer sqliteDB.Close()
	return true
}

func createDictionary(DBname string) {
	// # Create connection.
	sqliteDB, _ := sql.Open("sqlite3", DBname)
	defer sqliteDB.Close()
	// Dictionary prefab.
	sql := `
		CREATE TABLE IF NOT EXISTS Dictionary (
			"wordID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			"word" TEXT,
			"other" TEXT,
			"relShip" REAL
		);
	`
	// # Prepare before execution and cleanup.
	statement, err := sqliteDB.Prepare(sql)
	if err != nil {
		log.Fatal(err.Error())
	}

	statement.Exec()
	statement.Close()
}

func insertNewWord(DBName string, current string, succeeding string, score float32) {
	// # Create connection.
	sqliteDB, _ := sql.Open("sqlite3", DBName)
	defer sqliteDB.Close()
	sql := `
		INSERT INTO Dictionary (word, other, relShip)
			 VALUES (?,?,?)
	`
	// # Prepare before execution and cleanup.
	statement, err := sqliteDB.Prepare(sql)
	if err != nil {
		log.Fatal(err.Error())
	}

	statement.Exec(current, succeeding, score)
	statement.Close()

}

func checkWordPair(DBName, current, succeeding string) bool {
	// # Create connection.
	sqliteDB, _ := sql.Open("sqlite3", DBName)
	defer sqliteDB.Close()

	sql := `
		SELECT * FROM Dictionary
		 WHERE word = ?
		   AND other = ?
	`
	row, err := sqliteDB.Query(sql, current, succeeding)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer row.Close()
	result := false
	for row.Next() {
		result = true
	}

	return result
}

func incrementWordPair(DBName string, current string, succeeding string) {
	sqliteDB, _ := sql.Open("sqlite3", DBName)
	defer sqliteDB.Close()
	sql := `
		UPDATE Dictionary
		   SET relShip = relShip + 1
		 WHERE word = ?
		   AND other = ?
	`
	// # Prepare before execution and cleanup.
	statement, err := sqliteDB.Prepare(sql)
	if err != nil {
		log.Fatal(err.Error())
	}

	statement.Exec(current, succeeding)
	statement.Close()
}

func succeedingX(DBName, word string) ([]string, []float32) {
	sqliteDB, _ := sql.Open("sqlite3", DBName)
	defer sqliteDB.Close()
	sql := `
		SELECT other, relShip FROM Dictionary
		 WHERE word = ?
	`
	row, err := sqliteDB.Query(sql, word)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer row.Close()

	resOthers := make([]string, 0)
	resRelShips := make([]float32, 0)
	for row.Next() {
		var other string
		var relShip float32
		row.Scan(&other, &relShip)
		resOthers = append(resOthers, other)
		resRelShips = append(resRelShips, relShip)
	}

	return resOthers, resRelShips

}
