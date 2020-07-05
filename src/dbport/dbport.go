package dbport

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB(DBName string, overwriteOld bool) bool {
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

func modifier(DBName string, f func() (string, []interface{})) {
	sqliteDB, err := sql.Open("sqlite3", DBName)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer sqliteDB.Close()

	sql, bindings := f()

	statement, err := sqliteDB.Prepare(sql)

	if err != nil {
		log.Fatal(err.Error())
	}

	statement.Exec(bindings...)
	statement.Close()
}

func retriever(DBName string, f func() (string, []interface{})) ([]string, []float32) {
	sqliteDB, _ := sql.Open("sqlite3", DBName)
	defer sqliteDB.Close()

	sql, bindings := f()
	row, err := sqliteDB.Query(sql, bindings...)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer row.Close()

	// # Rigid impl.
	other, score := make([]string, 0, 100), make([]float32, 0, 100) // # 100 is arbitrary.
	for row.Next() {
		var otherTMP string
		var scoreTMP float32

		row.Scan(&otherTMP, &scoreTMP)

		other = append(other, otherTMP)
		score = append(score, scoreTMP)
	}
	return other, score
}

func CreateDictionary(DBName string) {
	modifier(DBName, func() (string, []interface{}) {
		sql := `
		CREATE TABLE IF NOT EXISTS Dictionary (
			"wordID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			"word" TEXT,
			"other" TEXT,
			"relShip" REAL
		);`

		return sql, make([]interface{}, 0)
	})
}

func InsertNewWord(DBName, current, succeeding string, score float32) {
	modifier(DBName, func() (string, []interface{}) {
		sql := `
		INSERT INTO Dictionary (word, other, relShip)
			 VALUES (?,?,?)
		`
		bindings := make([]interface{}, 3, 3)
		bindings[0], bindings[1], bindings[2] = current, succeeding, score
		return sql, bindings
	})
}

func IncrementWordPair(DBName, current, succeeding string) {
	modifier(DBName, func() (string, []interface{}) {
		sql := `
			UPDATE Dictionary
			SET relShip = relShip + 1
			WHERE word = ?
			AND other = ?
		`
		bindings := make([]interface{}, 2, 2)
		bindings[0], bindings[1] = current, succeeding
		return sql, bindings
	})
}

func CheckWordPair(DBName, current, other string) bool {
	others, _ := retriever(DBName, func() (string, []interface{}) {
		sql := `
			SELECT * FROM Dictionary
			WHERE word = ?
			AND other = ?
		`
		bindings := make([]interface{}, 2, 2)
		bindings[0], bindings[1] = current, other
		return sql, bindings
	})
	if len(others) > 0 {
		return true
	}
	return false
}

func SucceedingX(DBName, current string) ([]string, []float32) {
	return retriever(DBName, func() (string, []interface{}) {
		sql := `
			SELECT other, relShip FROM Dictionary
			WHERE word = ?
		`
		bindings := make([]interface{}, 1, 1)
		bindings[0] = current
		return sql, bindings
	})
}
