package sqlite

import (
	"database/sql"
	"log"

	"github.com/crunchypi/markov-go-sql.git/src/protocols"
	_ "github.com/mattn/go-sqlite3"
)

var _ protocols.DBAbstracter = (*SQLiteManager)(nil)

// SQLiteManager manages an sql connection
type SQLiteManager struct {
	db *sql.DB
}

// New creates a new sql connection and a template table
// (if one does not already exist).
func New(path string) (protocols.DBAbstracter, error) {
	man := SQLiteManager{}
	conn, err := sql.Open("sqlite3", path)
	man.db = conn

	man.createDictionary()
	return &man, err
}

// modifier is a wrapper for all functions which modify the db.
// argument should be a function which return an sql string and
// any bindings (as slice of interfaces)
func (s *SQLiteManager) modifier(query func() (string, []interface{})) {
	sql, bindings := query()
	statement, err := s.db.Prepare(sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec(bindings...)
	statement.Close()
}

// retriever is a wrapper for all functions which retrieve data
// from the db. Args:
// (1) query: function which returns an sql string and bindings
// (2) callback: called on each read, must take a ref to sql.Rows.
//   		   this is used to pull data from each row.
func (s *SQLiteManager) retriever(query func() (string, []interface{}),
	callback func(*sql.Rows)) {

	sql, bindings := query()
	row, err := s.db.Query(sql, bindings...)
	if err != nil {
		log.Fatal(err.Error())
	}

	for row.Next() {
		callback(row)
	}
}

// bindings does a common task in this file: converts
// arguments into a list of interfaces suitable for the sql pkg
func (s *SQLiteManager) bindings(word, other string, dst int) []interface{} {
	bindings := make([]interface{}, 3)
	bindings[0] = word
	bindings[1] = other
	bindings[2] = dst
	return bindings
}

// createDictionary creates a new template table (if not exists) in the db.
func (s *SQLiteManager) createDictionary() {
	s.modifier(func() (string, []interface{}) {
		sql := `
			CREATE TABLE IF NOT EXISTS Dictionary (
				"id" 	  	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"word" 	  	TEXT,
				"other"		TEXT,
				"distance"	INTEGER,
				"count" 	INTEGER
			);
		`
		return sql, make([]interface{}, 0)
	})
}

// ubsertNewPair inserts a new row describing a word relationship
// in the database. Count attribute is set to 1.
func (s *SQLiteManager) insertNewPair(word, other string, dst int) {
	s.modifier(func() (string, []interface{}) {
		sql := `
			INSERT INTO Dictionary (word,other,distance,count)
				 VALUES (?,?,?,1)
		`
		return sql, s.bindings(word, other, dst)
	})
}

// PairExists checks whether or not a row describing a word
// relationship exists in the database.
func (s *SQLiteManager) PairExists(word, other string, dst int) bool {
	f := func() (string, []interface{}) {
		sql := `
			SELECT word, other, distance FROM Dictionary
			 WHERE word = ?
			   AND other = ?
			   AND distance = ?
		`
		return sql, s.bindings(word, other, dst)
	}

	res := false
	callback := func(r *sql.Rows) {
		w, o, d := "", "", 0
		r.Scan(&w, &o, &d)
		res = w == word && o == other && d == dst
	}
	s.retriever(f, callback)
	return res
}

// IncrementPair updates a row describing a word relationship
// such that the count is incremented. Automatically creates a new
// pair with count = 1 if the pair does not already exist.
func (s *SQLiteManager) IncrementPair(word, other string, dst int) {
	if !s.PairExists(word, other, dst) {
		s.insertNewPair(word, other, dst)
		return
	}

	s.modifier(func() (string, []interface{}) {
		sql := `
			UPDATE Dictionary
			SET count = count + 1
				WHERE word = ?
				  AND other = ?
				  AND distance = ?
		`
		return sql, s.bindings(word, other, dst)
	})
}

// SucceedingX fetches all counterparts of the parameter such that
// the other word, distance and count can by anything. Returns
// a slice of type Record (defined in this file).
func (s *SQLiteManager) SucceedingX(word string) []protocols.WordRelationship {
	f := func() (string, []interface{}) {
		sql := `
			SELECT word, other, distance, count FROM Dictionary
			 WHERE word = ?
		`
		bindings := make([]interface{}, 1)
		bindings[0] = word
		return sql, bindings
	}

	records := make([]protocols.WordRelationship, 0, 100) // 100 is arbitrary
	callback := func(r *sql.Rows) {
		rec := protocols.WordRelationship{}
		r.Scan(&rec.Word, &rec.Other, &rec.Distance, &rec.Count)
		records = append(records, rec)
	}

	s.retriever(f, callback)
	return records
}
