package neo4j

import (
	"fmt"
	"strconv"

	"github.com/crunchypi/markov-go-sql.git/src/protocols"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// # Neo4jManager implements protocols.DBAbstracter
var _ protocols.DBAbstracter = (*Manager)(nil)

// Neo4jManager manages a neo4j connection and holds the set of
// methods required to implement protocols.DBAbstracter.
type Manager struct {
	db neo4j.Driver
}

// Used internal as args in main execution function.
type executeParams struct {
	cypher   string
	bindings map[string]interface{}
	callback func(neo4j.Result) // Optional
}

// New attempts to contact a Neo4j DB, given the params. Returns
// A Neo4jManager type (found in this file), which implements
// protocols.DBAbstracter.
func New(uri, user, pwd string, encr bool) (protocols.DBAbstracter, error) {
	new := Manager{}

	driver, err := neo4j.NewDriver(
		uri,
		neo4j.BasicAuth(user, pwd, ""),
		func(c *neo4j.Config) {
			c.Encrypted = encr
		})
	if err != nil {
		return &new, err
	}
	new.db = driver
	return &new, nil
}

// execute is the point of contact of the neo4j db/driver.
// Takes in 'executeParams' struct (found in this file),
// where 'callback' property (function) is optional.
func (n *Manager) execute(x executeParams) error {
	// # Open.
	session, err := n.db.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	// # Execute.
	res, err := session.Run(x.cypher, x.bindings)
	if err != nil {
		return err
	}

	// # Optional callback.
	if x.callback != nil {
		for res.Next() {
			x.callback(res)
		}
	}
	return nil
}

// newNodes attempts to create a node for each element in 'words'.
// Done with MERGE instead of CREATE (neo4j syntax) for an
// all-or-nothing execution.
func (n *Manager) newNodes(words []string) error {
	cypher := ``
	bindings := make(map[string]interface{})
	for i, v := range words {
		cypher += fmt.Sprintf(`MERGE (:MChain{word:$%d})`, i)
		bindings[strconv.Itoa(i)] = v
	}

	return n.execute(executeParams{
		cypher:   cypher,
		bindings: bindings,
	})
}

// IncrementPair attempts to increment the 'count' property on
// the relationship between 'word' and 'other' with a certain 'dst'.
// If this is not possible, a new relationship will be created.
// Note: 'word' and 'other' do not have to be in the db.
func (n *Manager) IncrementPair(word, other string, dst int) {
	// # Add nodes if they do not already exist - makes cypher simpler.
	n.newNodes([]string{word, other})
	cypher := `
		// If relationship exists, increment it
		OPTIONAL MATCH (a:MChain{word:$word})-[b:conn{dst:$dst}]->(c:MChain{word:$other})
		SET b.count = b.count + 1

		// Else - create relship only if nodes exist.
		WITH a AS _
		MATCH (x:MChain{word:$word}), (y:MChain{word:$other})
		WHERE NOT (x)-[:conn]->(y)
		CREATE (x)-[z:conn{dst:$dst, count:1}]->(y)
	`
	bindings := map[string]interface{}{
		"word":  word,
		"other": other,
		"dst":   dst,
	}
	n.execute(executeParams{
		cypher:   cypher,
		bindings: bindings,
	})
}

func (n *Manager) SucceedingX(word string) []protocols.WordRelationship {

	res := make([]protocols.WordRelationship, 0, 100) // # 100 is arbitrary
	n.execute(executeParams{
		cypher: `
			 MATCH (x:MChain{word:$word})-[r:conn]->(y)
			RETURN y.word AS a, x.word AS b, r.dst AS c, r.count AS d 
		`,
		bindings: map[string]interface{}{
			"word": word,
		},
		callback: func(r neo4j.Result) {
			a, aok := r.Record().Get("a")
			b, bok := r.Record().Get("b")
			c, cok := r.Record().Get("c")
			d, dok := r.Record().Get("d")
			// log.Printf("####, %s, %s, %d, %d", a, b, c, d)
			if aok && bok && cok && dok {
				newNode := protocols.WordRelationship{
					Word:     a.(string),
					Other:    b.(string),
					Distance: int(c.(int64)),
					Count:    int(d.(int64)),
				}
				res = append(res, newNode)

			}
		},
	})
	return res
}

// If relationship exists, increment it
// OPTIONAL MATCH (a:MChain{word:"a"})-[b:conn]->(c:MChain{word:"b"})
// FOREACH ( _ in CASE WHEN a.word="a" THEN [1] ELSE [0] END |
//   SET b.dst = b.dst + 1)

// // Create relship only if nodes exist
// WITH a AS _
// MATCH (x:MChain{word:"a"}), (y:MChain{word:"b"})
// WHERE NOT (x)-[:conn]->(y)
// CREATE (x)-[z:conn{dst:1}]->(y)
