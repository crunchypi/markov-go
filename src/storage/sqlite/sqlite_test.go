package sqlite

import (
	"os"
	"testing"
)

var path = "test.sqlite"

func cleanup() {
	os.Remove(path)
}

// Have to check manually for this.
func TestNew(t *testing.T) {
	_, err := New(path)
	// defer cleanup()

	if err != nil {
		t.Error("Failed on db creation", err)
	}
}

// Have to check manually for this.
func TestInsertNewPair(t *testing.T) {
	da, _ := New(path)
	db := da.(*SQLiteManager)
	// defer cleanup()
	db.insertNewPair("one", "two", 3)
}

func TestPairExists(t *testing.T) {
	da, _ := New(path)
	db := da.(*SQLiteManager)
	defer cleanup()

	word, other, dst := "one", "two", 3

	db.insertNewPair(word, other, dst)
	exists, _ := db.PairExists(word, other, dst)
	if !exists {
		t.Error("issue with .PairExists")
	}
}

// Have to check manually for this.
func TestIncrementPair(t *testing.T) {
	da, _ := New(path)
	db := da.(*SQLiteManager)
	// defer cleanup()

	word, other, dst := "one", "two", 3
	db.IncrementPair(word, other, dst)
	db.IncrementPair(word, other, dst)
	// # Check DB, this relationship should have 2 count.
}

func TestSucceedingX(t *testing.T) {
	da, _ := New(path)
	db := da.(*SQLiteManager)
	defer cleanup()

	wild := "wild"
	word1, other1, dst1 := "one_1", "two_1", 1
	word2, other2, dst2 := "one_2", "two_2", 2

	db.IncrementPair(word1, other1, dst1)
	db.IncrementPair(word1, other1, dst1)
	db.IncrementPair(word1, wild, dst1)
	db.IncrementPair(word2, other2, dst2)

	res1, _ := db.SucceedingX(word1)
	res2, _ := db.SucceedingX(word2)

	if len(res1) != 2 || len(res2) != 1 {
		t.Error("failed to fetch")
	}

	if res1[0].Word != word1 || res1[0].Other != other1 ||
		res1[0].Distance != dst1 || res1[0].Count != 2 {
		t.Error("incorrect result for class 1")
	}
	if res1[1].Word != word1 || res1[1].Other != wild ||
		res1[1].Distance != dst1 || res1[1].Count != 1 {
		t.Error("incorrect result for class 1")
	}
	if res2[0].Word != word2 || res2[0].Other != other2 ||
		res2[0].Distance != dst2 || res2[0].Count != 1 {
		t.Error("incorrect result for class 2")
	}

}
