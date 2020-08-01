package neo4j

import (
	"fmt"
	"strconv"

	"github.com/crunchypi/markov-go-sql.git/src/protocols"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// # Neo4jManager implements protocols.DBAbstracter
var _ protocols.DBAbstracter = (*Neo4jManager)(nil)

type Neo4jManager struct {
	db neo4j.Driver
}

// Used internal as args in main execution function.
type executeParams struct {
	cypher   string
	bindings map[string]interface{}
	callback func(neo4j.Result) // Optional
}

func New(uri, user, pwd string, encr bool) (protocols.DBAbstracter, error) {
	new := Neo4jManager{}

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

// modifier is the point of contact of the neo4j db/driver.
// @ TODO: transactions + batched commands?
func (n *Neo4jManager) execute(x executeParams) error {
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

// bindings does a common task in this file: converts
// arguments into a map suitable for the sql pkg
func (n *Neo4jManager) bindings(word, other string, dst int) map[string]interface{} {
	return map[string]interface{}{
		"word":  word,
		"other": other,
		"dst":   dst,
	}
}

func (n *Neo4jManager) newNodes(words []string) error {
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

func (n *Neo4jManager) IncrementPair(word, other string, dst int) {
	n.newNodes([]string{word, other})
	cypher := `
		// If relationship exists, increment it
		OPTIONAL MATCH (a:MChain{word:$word})-[b:conn{dst:$dst}]->(c:MChain{word:$other})
		SET b.count = b.count + 1

		// Create relship only if nodes exist
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

func (n *Neo4jManager) SucceedingX(word string) []protocols.WordRelationship {

	res := make([]protocols.WordRelationship, 0, 10) // # 10 is arbitrary
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
