package dbport

import "testing"

func TestCreateDB(t *testing.T) {
	// # Results meant to be checked manually in SQLIte.
	DBName := "test.sqlite"
	CreateDB(DBName, true)
	CreateDictionary(DBName)
}

func TestNewWordPair(t *testing.T) {
	// # Results meant to be checked manually in SQLIte.
	DBName := "test.sqlite"
	CreateDB(DBName, true)
	CreateDictionary(DBName)
	InsertNewWord(DBName, "one", "two", 3.)
}

func TestIncrementWordPair(t *testing.T) {
	// # Results meant to be checked manually in SQLIte.
	DBName, current, other := "test.sqlite", "one", "two"
	CreateDB(DBName, true)
	CreateDictionary(DBName)
	InsertNewWord(DBName, current, other, 0)
	IncrementWordPair(DBName, current, other)
}

func TestCheckWordPair(t *testing.T) {
	DBName, current, other := "test.sqlite", "one", "two"
	CreateDB(DBName, true)
	CreateDictionary(DBName)
	InsertNewWord(DBName, current, other, 0)
	res := CheckWordPair(DBName, current, other)
	if !res {
		t.Error("created word pair but coult not verify it in the DB.")
	}
}

func TestSucceedingX(t *testing.T) {
	DBName, current, other := "test.sqlite", "one", "two"
	CreateDB(DBName, true)
	CreateDictionary(DBName)
	InsertNewWord(DBName, current, other, 0)
	res, _ := SucceedingX(DBName, current)
	if res[0] != other {
		t.Error("created word pair but could not retrieve it.", res)
	}
}
